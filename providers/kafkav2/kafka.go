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
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	writer "github.com/erda-project/erda-infra/pkg/parallel-writer"
)

// Interface .
type Interface interface {
	NewProducer(c ProducerConfig, options ...ProducerOption) (writer.Writer, error)
	Servers() string
}

// Producer .
type Producer interface {
	writer.Writer
}

type config struct {
	Servers string `file:"servers" env:"BOOTSTRAP_SERVERS" default:"localhost:9092" desc:"kafka servers"`
}

// provider .
type provider struct {
	Cfg *config
	Log logs.Logger
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
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

func init() {
	servicehub.Register("kafka-v2", &servicehub.Spec{
		Services: []string{"kafka-v2", "kafka-producer-v2"},
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
		},
		ConfigFunc: func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
