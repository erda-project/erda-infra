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
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

type config struct {
	Message string `file:"message" flag:"msg" default:"hi" desc:"message to show" env:"HELLO_MESSAGE"`
}

type provider struct {
	Cfg     *config
	Log     logs.Logger
	closeCh chan struct{}
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.Log.Info("message: ", p.Cfg.Message)
	return nil
}

func (p *provider) Start() error {
	p.Log.Info("hello provider is starting...")
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			p.Log.Info("do something...")
		case <-p.closeCh:
			return nil
		}
	}
}

func (p *provider) Close() error {
	p.Log.Info("hello provider is closing...")
	close(p.closeCh)
	return nil
}

func init() {
	servicehub.Register("hello-provider", &servicehub.Spec{
		Services:    []string{"hello"},
		Description: "hello for example",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{
				closeCh: make(chan struct{}),
			}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
