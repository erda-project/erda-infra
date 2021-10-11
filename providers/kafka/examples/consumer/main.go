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

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/kafka"
)

type config struct {
	Input kafka.ConsumerConfig `file:"input"`
}

type provider struct {
	Cfg   *config
	Log   logs.Logger
	Kafka kafka.Interface `autowired:"kafka-consumer"`
}

func (p *provider) Run(ctx context.Context) error {
	p.Kafka.NewConsumer(&p.Cfg.Input, p.invoke)
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *provider) invoke(key []byte, value []byte, topic *string, timestamp time.Time) error {
	fmt.Println(string(value))
	return nil
}

func init() {
	servicehub.Register("examples", &servicehub.Spec{
		Services:   []string{"hello"},
		ConfigFunc: func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
