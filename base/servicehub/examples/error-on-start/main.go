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
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

type config struct {
	Error bool `file:"error"`
}

type provider struct {
	Cfg *config
	Log logs.Logger
}

func (p *provider) Run(ctx context.Context) error {
	if p.Cfg.Error {
		time.Sleep(3 * time.Second)
		return fmt.Errorf("run error")
	}
	p.Log.Info("run with no error")
	for {
		select {
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
	servicehub.Run(&servicehub.RunOptions{
		Content: `
hello-provider:
    error: false
hello-provider@error1:
    error: true
hello-provider@error2:
    error: true
`,
	})
}
