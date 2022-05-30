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
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// BatchReader .
type BatchReader interface {
	ReadN(buf []interface{}, timeout time.Duration) (int, error)
	Confirm() error
	Close() error
}

// BatchReaderConfig .
type BatchReaderConfig struct {
	Topics  []string               `file:"topics" desc:"topics"`
	Group   string                 `file:"group" desc:"consumer group id"`
	Options map[string]interface{} `file:"options" desc:"options"`
}

// BatchReaderOption .
type BatchReaderOption interface{}

// WithReaderDecoder .
func WithReaderDecoder(dec Decoder) BatchReaderOption {
	return BatchReaderOption(dec)
}

// Decoder .
type Decoder func(key, value []byte, topic *string, timestamp time.Time) (interface{}, error)

func (s *service) NewBatchReader(cfg *BatchReaderConfig, options ...BatchReaderOption) (BatchReader, error) {
	var dec Decoder
	for _, opt := range options {
		switch v := opt.(type) {
		case Decoder:
			dec = v
		}
	}
	if dec == nil {
		dec = func(key, value []byte, topic *string, timestamp time.Time) (interface{}, error) {
			return value, nil
		}
	}
	kc := convertToConfigMap(mergeMap(s.p.Cfg.Comsumer.Options, cfg.Options))
	return newKafkaReader(s.p.Cfg.Servers, cfg.Group, cfg.Topics, kc, dec)
}

func newKafkaReader(servers, group string, topics []string, kc kafka.ConfigMap, dec Decoder) (BatchReader, error) {
	kc["bootstrap.servers"] = servers
	kc["group.id"] = group
	kc["enable.auto.offset.store"] = false
	kc["enable.auto.commit"] = false
	delete(kc, "auto.offset.reset")
	delete(kc, "auto.commit.interval.ms")
	return &kafkaBatchReader{
		kc:     kc,
		topics: topics,
		decode: dec,
	}, nil
}

type kafkaBatchReader struct {
	kc       kafka.ConfigMap
	topics   []string
	consumer *kafka.Consumer
	decode   Decoder
}

func (r *kafkaBatchReader) ReadN(buf []interface{}, timeout time.Duration) (int, error) {
	if r.consumer == nil {
		consumer, err := kafka.NewConsumer(&r.kc)
		if err != nil {
			return 0, err
		}
		err = consumer.SubscribeTopics(r.topics, nil)
		if err != nil {
			consumer.Close()
			return 0, err
		}
		r.consumer = consumer
	}
	size := len(buf)
	var offset int
	maxWaitTimer := time.NewTimer(timeout * 3)
	defer maxWaitTimer.Stop()
	for {
		if offset >= size {
			break
		}
		select {
		case <-maxWaitTimer.C:
			break
		default:
		}
		msg, err := r.consumer.ReadMessage(timeout)
		if err != nil {
			if kerr, ok := err.(kafka.Error); ok {
				if kerr.Code() == kafka.ErrTimedOut {
					return offset, nil
				}
			}
			r.Close()
			return offset, err
		}

		data, err := r.decode(msg.Key, msg.Value, msg.TopicPartition.Topic, msg.Timestamp)
		if err != nil {
			// ingore decode error
			continue
		}

		_, err = r.consumer.StoreOffsets([]kafka.TopicPartition{msg.TopicPartition})
		if err != nil {
			return offset, err
		}
		buf[offset] = data
		offset++
	}
	return offset, nil
}

func (r *kafkaBatchReader) Confirm() error {
	if r.consumer != nil {
		_, err := r.consumer.Commit()
		if kerr, ok := err.(kafka.Error); ok {
			if kerr.Code() == kafka.ErrNoOffset {
				return nil
			}
		}
		return err
	}

	return nil
}

func (r *kafkaBatchReader) Close() error {
	consumer := r.consumer
	if consumer != nil {
		consumer.Unsubscribe()
		err := consumer.Close()
		r.consumer = nil
		return err
	}
	return nil
}

// CommittedOffsets .
func (r *kafkaBatchReader) CommittedOffsets() ([]kafka.TopicPartition, error) {
	consumer := r.consumer
	if consumer == nil {
		return nil, nil
	}
	ps, err := consumer.Assignment()
	if err != nil {
		return nil, err
	}
	return consumer.Committed(ps, 30*1000)
}

// CommittedOffsets .
func CommittedOffsets(r BatchReader) ([]kafka.TopicPartition, error) {
	bw, ok := r.(*kafkaBatchReader)
	if ok {
		return bw.CommittedOffsets()
	}
	return nil, nil
}
