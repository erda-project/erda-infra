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

package caller

import (
	"context"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/examples/service/protocol/pb"
	"github.com/erda-project/erda-infra/pkg/transport"
)

type config struct {
	Name string `file:"name" default:"recallsong"`
}

// +provider
type provider struct {
	Cfg     *config
	Log     logs.Logger
	Greeter pb.GreeterServiceServer // remote or local service. this doesn't need to care
}

// Run this is an optional
func (p *provider) Run(ctx context.Context) error {
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			header := transport.Header{}
			header.Set("Custom-Header", "Custom-Header-Value")
			resp, err := p.Greeter.SayHello(transport.WithHeader(context.Background(), header), &pb.HelloRequest{
				Name: p.Cfg.Name,
			})
			if err != nil {
				p.Log.Error(err)
			}
			p.Log.Info(resp)
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.Register("caller", &servicehub.Spec{
		Services:     []string{},
		Description:  "this is caller example",
		Dependencies: []string{"erda.infra.example.GreeterService"},
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
