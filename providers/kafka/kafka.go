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
	"reflect"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	writer "github.com/erda-project/erda-infra/pkg/parallel-writer"
)

// Interface .
type Interface interface {
	NewConsumer(c *ConsumerConfig, handler ConsumerFunc, options ...ConsumerOption) error
	NewConsumerWitchCreator(c *ConsumerConfig, creator func(i int) (ConsumerFunc, error), options ...ConsumerOption) error
	NewProducer(c *ProducerConfig, options ...ProducerOption) (writer.Writer, error)
	Servers() string
	ProduceChannelSize() int
	ProduceEventsChannelSize() int
	NewAdminClient() (*kafka.AdminClient, error)
}

// Producer .
type Producer interface {
	writer.Writer
	ProduceChannelSize() int
	EventsChannelSize() int
}

type config struct {
	Servers  string `file:"servers" env:"BOOTSTRAP_SERVERS" default:"localhost:9092" desc:"kafka servers"`
	Producer struct {
		Options map[string]interface{} `file:"options"`
	} `file:"producer"`
	Comsumer struct {
		Options map[string]interface{} `file:"options"`
	} `file:"comsumer"`
	Admin struct {
		Options map[string]interface{} `file:"options"`
	} `file:"admin"`
}

// provider .
type provider struct {
	Cfg      *config
	Log      logs.Logger
	producer sharedProducer
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.producer.log = p.Log
	return nil
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, options ...interface{}) interface{} {
	return &service{
		p:    p,
		log:  p.Log.Sub(ctx.Caller()),
		name: ctx.Caller(),
	}
}

type service struct {
	p    *provider
	log  logs.Logger
	name string
}

var _ Interface = (*service)(nil)

func (s *service) Servers() string { return s.p.Cfg.Servers }

func (s *service) NewAdminClient() (*kafka.AdminClient, error) {
	kc := kafka.ConfigMap{}
	if s.p.Cfg.Admin.Options != nil {
		for k, v := range s.p.Cfg.Admin.Options {
			kc[k] = v
		}
	}
	kc["bootstrap.servers"] = s.p.Cfg.Servers
	return kafka.NewAdminClient(&kc)
}

func init() {
	servicehub.Register("kafka", &servicehub.Spec{
		Services: []string{"kafka", "kafka-producer", "kafka-consumer"},
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
		},
		Description: "kafka、kafka-producer、kafka-consumer",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
