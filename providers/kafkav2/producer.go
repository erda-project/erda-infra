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

package kafkav2

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/recallsong/go-utils/reflectx"
	"github.com/segmentio/kafka-go"

	"github.com/erda-project/erda-infra/base/logs"
	writer "github.com/erda-project/erda-infra/pkg/parallel-writer"
)

// Message .
type Message struct {
	Topic string
	Data  []byte
	Key   []byte
}

// ProducerConfig .
type ProducerConfig struct {
	Topic       string        `file:"topic"`
	Parallelism uint64        `file:"parallelism"  default:"3" env:"PROVIDER_KAFKA_V2_PRODUCER_PARALLELISM"`
	Async       bool          `file:"async" default:"true" env:"PROVIDER_KAFKA_V2_PRODUCER_ASYNC"`
	Timeout     time.Duration `file:"timeout" default:"30s" env:"PROVIDER_KAFKA_V2_PRODUCER_TIMEOUT"`
	Batch       struct {
		Size      int           `file:"size" default:"100" env:"PROVIDER_KAFKA_V2_PRODUCER_BATCH_SIZE"`
		SizeBytes int64         `file:"size_bytes" default:"1048576" env:"PROVIDER_KAFKA_V2_PRODUCER_BATCH_SIZE_BYTES"`
		Timeout   time.Duration `file:"timeout" default:"800ms" env:"PROVIDER_KAFKA_V2_PRODUCER_BATCH_TIMEOUT"`
	} `file:"batch"`
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

func newProducer(servers string, cfg ProducerConfig, log logs.Logger) (*producer, error) {
	prod := &producer{
		logger: log,
	}
	pw := &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(servers, ",")...),
		Balancer:               kafka.CRC32Balancer{},
		Async:                  cfg.Async,
		AllowAutoTopicCreation: true,
		WriteTimeout:           cfg.Timeout,
		BatchSize:              cfg.Batch.Size,
		BatchTimeout:           cfg.Batch.Timeout,
		BatchBytes:             cfg.Batch.SizeBytes,
	}
	if cfg.Topic != "" {
		pw.Topic = cfg.Topic
	}
	prod.pw = pw
	return prod, nil
}

type producer struct {
	logger logs.Logger
	pw     *kafka.Writer
}

func (p *producer) Write(data interface{}) error {
	return p.publish(data)
}

func (p *producer) WriteN(data ...interface{}) (int, error) {
	for i, item := range data {
		err := p.publish(item)
		if err != nil {
			return i, err
		}
	}
	return len(data), nil
}

func (p *producer) publish(data interface{}) error {
	var (
		value []byte
		key   []byte
	)
	topic := ""

	switch val := data.(type) {
	case Message:
		if val.Topic != "" {
			topic = val.Topic
		}
		value = val.Data
		key = val.Key
	case []byte:
		value = val
	case string:
		value = reflectx.StringToBytes(val)
	default:
		data, err := json.Marshal(data)
		if err != nil {
			return err
		}
		value = data
	}
	if p.pw.Topic == "" {
		err := p.pw.WriteMessages(context.TODO(), kafka.Message{
			Topic: topic,
			Key:   key,
			Value: value,
		})
		if err != nil {
			return err
		}
	} else {
		err := p.pw.WriteMessages(context.TODO(), kafka.Message{
			Key:   key,
			Value: value,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *producer) Close() error {
	return p.pw.Close()
}

func (s *service) NewProducer(cfg ProducerConfig, options ...ProducerOption) (writer.Writer, error) {
	var eh writer.ErrorHandler = s.producerError
	for _, item := range options {
		if item != nil && item.errHandler() != nil {
			eh = item.errHandler()
		}
	}
	prod, err := newProducer(s.p.Cfg.Servers, cfg, s.log)
	if err != nil {
		return nil, err
	}
	return writer.ParallelBatch(func(uint64) writer.Writer {
		return prod
	}, cfg.Parallelism, 1, 0, eh), nil
}

func (s *service) producerError(err error) error {
	s.log.Errorf("fail to write kafka: %s", err)
	return nil // skip error
}
