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
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

type config struct {
	Message string `file:"message" flag:"msg" default:"hi" desc:"message to show" env:"HELLO_MESSAGE"`
}

type provider struct {
	Cfg *config
	Log logs.Logger
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.Log.Info("message: ", p.Cfg.Message)
	ctx.AddTask(p.Task1)
	ctx.AddTask(p.Task2, servicehub.WithTaskName("task2"))
	return nil
}

func (p *provider) Task1(ctx context.Context) error {
	p.Log.Info("Task1 is running...")
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			p.Log.Info("Task1 do something...")
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *provider) Task2(ctx context.Context) error {
	p.Log.Info("Task2 is running...")
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			p.Log.Info("Task2 do something...")
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *provider) Run(ctx context.Context) error {
	p.Log.Info("provider is running...")
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			p.Log.Info("Run do something...")
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.Register("hello-provider", &servicehub.Spec{
		Services:    []string{"hello"},
		Description: "hello for example",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}

// INFO[2021-08-04 16:33:10.977] message: hello world                          module=hello-provider
// INFO[2021-08-04 16:33:10.977] provider hello-provider initialized
// INFO[2021-08-04 16:33:10.977] signals to quit: [hangup interrupt terminated quit]
// INFO[2021-08-04 16:33:10.977] provider hello-provider task(task2) running ...
// INFO[2021-08-04 16:33:10.978] Task2 is running...                           module=hello-provider
// INFO[2021-08-04 16:33:10.978] provider hello-provider task(1) running ...
// INFO[2021-08-04 16:33:10.978] Task1 is running...                           module=hello-provider
// INFO[2021-08-04 16:33:10.978] provider hello-provider running ...
// INFO[2021-08-04 16:33:10.978] provider is running...                        module=hello-provider
// INFO[2021-08-04 16:33:13.980] Task1 do something...                         module=hello-provider
// INFO[2021-08-04 16:33:13.980] Task2 do something...                         module=hello-provider
// INFO[2021-08-04 16:33:13.983] Run do something...                           module=hello-provider
// INFO[2021-08-04 16:33:16.982] Run do something...                           module=hello-provider
// INFO[2021-08-04 16:33:16.982] Task2 do something...                         module=hello-provider
// INFO[2021-08-04 16:33:16.982] Task1 do something...                         module=hello-provider
// INFO[2021-08-04 16:33:19.981] Run do something...                           module=hello-provider
// INFO[2021-08-04 16:33:19.981] Task2 do something...                         module=hello-provider
// INFO[2021-08-04 16:33:19.981] Task1 do something...                         module=hello-provider
// ^C
// INFO[2021-08-04 16:33:20.247] provider hello-provider Run exit
// INFO[2021-08-04 16:33:20.247] provider hello-provider task(task2) exit
// INFO[2021-08-04 16:33:20.247] provider hello-provider task(1) exit
