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

// ConsumerConfig .
type ConsumerConfig struct {
	Topics      []string               `file:"topics" desc:"topics"`
	Group       string                 `file:"group" desc:"consumer group id"`
	Parallelism uint64                 `file:"parallelism" desc:"parallelism"`
	Options     map[string]interface{} `file:"options" desc:"options"`
}

// ConsumerFunc .
type ConsumerFunc func(key []byte, value []byte, topic *string, timestamp time.Time) error

func mergeMap(a, b map[string]interface{}) map[string]interface{} {
	if a == nil || len(a) == 0 {
		return b
	}
	if b == nil || len(b) == 0 {
		return a
	}
	c := make(map[string]interface{}, len(a)+len(b))
	for k, v := range a {
		c[k] = v
	}
	for k, v := range b {
		c[k] = v
	}
	return c
}

func (s *service) NewConsumer(cfg *ConsumerConfig, handler ConsumerFunc, options ...ConsumerOption) error {
	return s.NewConsumerWitchCreator(cfg, func(int) (ConsumerFunc, error) { return handler, nil })
}

func (s *service) NewConsumerWitchCreator(cfg *ConsumerConfig, creator func(i int) (ConsumerFunc, error), opts ...ConsumerOption) error {
	options := mergeMap(s.p.Cfg.Comsumer.Options, cfg.Options)
	parallelism := int(cfg.Parallelism)
	var consumerListener func(i int, c *kafka.Consumer)
	for _, opt := range opts {
		switch v := opt.(type) {
		case func(i int, c *kafka.Consumer):
			consumerListener = v
		}
	}
	for i := 0; i < parallelism; i++ {
		var kc kafka.ConfigMap
		if options != nil {
			kc = convertToConfigMap(options)
		} else {
			kc = kafka.ConfigMap{}
		}
		kc["bootstrap.servers"] = s.p.Cfg.Servers
		kc["group.id"] = cfg.Group

		handler, err := creator(i)
		if err != nil {
			return err
		}
		go func(i int, handler ConsumerFunc) {
		loop:
			for {
				consumer, err := kafka.NewConsumer(&kc)
				if err != nil {
					s.log.Errorf("failed to create kafka consumer with config %#v: %v", cfg, err)
					time.Sleep(3 * time.Second)
					continue
				}
				if err = consumer.SubscribeTopics(cfg.Topics, nil); err != nil {
					s.log.Errorf("failed to subscribe kafka topics %v: %#v", cfg.Topics, err)
					s.logError(consumer.Close())
					time.Sleep(3 * time.Second)
					continue
				}
				if consumerListener != nil {
					consumerListener(i, consumer)
				}
				s.log.Infof("create kafka consumer %d with topics: %v, config: %#v", i, cfg.Topics, kc)
				var errors int
				for {
					message, err := consumer.ReadMessage(-1)
					if err != nil {
						var topic string
						if message != nil {
							topic = *message.TopicPartition.Topic
						}
						s.log.Errorf("topic: %s .fail to read message from kafka: %v", topic, err)
						errors++
						if errors > 10 {
							s.logError(consumer.Close())
							time.Sleep(1 * time.Second)
							continue loop
						}
						continue
					}
					err = handler(message.Key, message.Value, message.TopicPartition.Topic, message.Timestamp)
					if err != nil {
						s.log.Errorf("fail to process message: %v", err)
					}
					errors = 0
				}
			}
		}(i, handler)
	}
	return nil
}

func convertToConfigMap(m map[string]interface{}) kafka.ConfigMap {
	cm := make(kafka.ConfigMap, len(m))
	for k, v := range m {
		cm[k] = v
	}
	return cm
}

func (s *service) logError(err error) {
	if err != nil {
		s.log.Error(err)
	}
}

// ConsumerOption .
type ConsumerOption interface{}

// WithConsumerListener .
func WithConsumerListener(fn func(i int, c *kafka.Consumer)) ConsumerOption {
	return ConsumerOption(fn)
}
