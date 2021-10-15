// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/recallsong/go-utils/reflectx"

	"github.com/erda-project/erda-infra/base/logs"
	writer "github.com/erda-project/erda-infra/pkg/parallel-writer"
)

// Message .
type Message struct {
	Topic *string
	Data  []byte
	Key   []byte
}

// StringMessage .
type StringMessage struct {
	Topic *string
	Data  string
}

// ProducerConfig .
type ProducerConfig struct {
	Topic       string `file:"topic" env:"KAFKA_P_TOPIC" desc:"topic"`
	Parallelism uint64 `file:"parallelism" env:"KAFKA_P_PARALLELISM" default:"4" desc:"parallelism"`
	Batch       struct {
		Size    uint64        `file:"size" env:"KAFKA_P_BATCH_SIZE" default:"100" desc:"batch size"`
		Timeout time.Duration `file:"timeout" env:"KAFKA_P_BUFFER_TIMEOUT" default:"30s" desc:"timeout to flush buffer for batch write"`
	} `file:"batch"`
	Shared  bool                   `file:"shared" default:"true" desc:"shared producer instance"`
	Options map[string]interface{} `file:"options" desc:"options"`
}

// ProducerOption .
type ProducerOption interface {
	errHandler() func(error) error
}

type producerOption struct{ _eh func(error) error }

func (p *producerOption) errHandler() func(error) error { return p._eh }

// WithAsyncWriteErrorHandler .
func WithAsyncWriteErrorHandler(eh func(error) error) ProducerOption {
	return &producerOption{_eh: eh}
}

func newProducer(servers string, extra map[string]interface{}, log logs.Logger) (*kafka.Producer, error) {
	kc := kafka.ConfigMap{"go.batch.producer": true}
	if extra != nil {
		for k, v := range extra {
			kc[k] = v
		}
	}
	kc["bootstrap.servers"] = servers
	kp, err := kafka.NewProducer(&kc)
	if err != nil {
		return nil, err
	}
	go consumeEvents(kp, log)
	return kp, err
}

func consumeEvents(kp *kafka.Producer, log logs.Logger) {
	for e := range kp.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Errorf("Kafka delivery failed: %v", ev.TopicPartition)
			}
		}
	}
	log.Debugf("exit kafka events consumer")
}

type sharedProducer struct {
	lock     sync.Mutex
	instance *kafka.Producer
	refs     int
	log      logs.Logger
}

func (p *sharedProducer) release() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.refs == 0 {
		return nil
	}
	p.refs--
	if p.refs == 0 {
		p.instance.Close()
	}
	return nil
}

func (p *sharedProducer) get(servers string, extra map[string]interface{}) (*kafka.Producer, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.refs == 0 {
		kp, err := newProducer(servers, extra, p.log)
		if err != nil {
			return nil, err
		}
		p.instance = kp
	}
	p.refs++
	return p.instance, nil
}

type producer struct {
	kp    *kafka.Producer
	close func() error
	topic string
}

func (p *producer) ProduceChannelSize() int {
	return len(p.kp.ProduceChannel())
}

func (p *producer) EventsChannelSize() int {
	return len(p.kp.Events())
}

func (p *producer) Write(data interface{}) error {
	delivery := make(chan kafka.Event)
	if err := p.publish(data, delivery); err != nil {
		return err
	}
	select {
	case <-delivery:
		return nil
	default:
		// wait 1s
		p.kp.Flush(1000)
	}
	return nil
}

func (p *producer) WriteN(data ...interface{}) (int, error) {
	delivery := make(chan kafka.Event)
	for i, item := range data {
		err := p.publish(item, delivery)
		if err != nil {
			return i, err
		}
	}
	for i := 0; i < len(data); {
		select {
		case <-delivery:
			i++
		default:
			// wait 1s
			p.kp.Flush(1000)
		}
	}
	return len(data), nil
}

func (p *producer) publish(data interface{}, delivery chan kafka.Event) error {
	var (
		bytes []byte
		key   []byte
	)
	topic := &p.topic
	switch val := data.(type) {
	case *Message:
		if val.Topic != nil {
			topic = val.Topic
		}
		bytes = val.Data
		key = val.Key
	case *StringMessage:
		if val.Topic != nil {
			topic = val.Topic
		}
		bytes = reflectx.StringToBytes(val.Data)
	case []byte:
		bytes = val
	case string:
		bytes = reflectx.StringToBytes(val)
	default:
		data, err := json.Marshal(data)
		if err != nil {
			return err
		}
		bytes = data
	}
	return p.kp.Produce(&kafka.Message{
		Value:          bytes,
		Key:            key,
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
	}, delivery)
}

func (p *producer) Close() error {
	return p.close()
}

func (s *service) NewProducer(c *ProducerConfig, options ...ProducerOption) (writer.Writer, error) {
	var eh writer.ErrorHandler = s.producerError
	for _, item := range options {
		if item != nil && item.errHandler() != nil {
			eh = item.errHandler()
		}
	}
	if c.Shared {
		kp, err := s.p.producer.get(s.p.Cfg.Servers, s.p.Cfg.Producer.Options)
		if err != nil {
			return nil, err
		}
		return writer.ParallelBatch(func(uint64) writer.Writer {
			return &producer{
				kp:    kp,
				close: s.p.producer.release,
				topic: c.Topic,
			}
		}, c.Parallelism, c.Batch.Size, c.Batch.Timeout, eh), nil
	}
	kp, err := newProducer(s.p.Cfg.Servers, c.Options, s.log)
	if err != nil {
		return nil, err
	}
	return writer.ParallelBatch(func(uint64) writer.Writer {
		return &producer{
			kp: kp,
			close: func() error {
				kp.Close()
				return nil
			},
			topic: c.Topic,
		}
	}, c.Parallelism, c.Batch.Size, c.Batch.Timeout, eh), nil
}

func (s *service) producerError(err error) error {
	s.log.Errorf("fail to write kafka: %s", err)
	return nil // skip error
}

func (s *service) ProduceChannelSize() int {
	s.p.producer.lock.Lock()
	defer s.p.producer.lock.Unlock()
	if s.p.producer.instance == nil {
		return 0
	}
	return len(s.p.producer.instance.ProduceChannel())
}

func (s *service) ProduceEventsChannelSize() int {
	s.p.producer.lock.Lock()
	defer s.p.producer.lock.Unlock()
	if s.p.producer.instance == nil {
		return 0
	}
	return len(s.p.producer.instance.Events())
}
