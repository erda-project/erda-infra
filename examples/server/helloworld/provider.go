// Copyright 2021 Terminus
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

package example

import (
	"context"
	"fmt"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/examples/protocol/pb"
	"github.com/erda-project/erda-infra/pkg/transport"
	"github.com/erda-project/erda-infra/pkg/transport/interceptor"
)

type config struct {
}

type provider struct {
	Cfg      *config
	Log      logs.Logger
	Register transport.Register
}

func (p *provider) Init(ctx servicehub.Context) error {
	// TODO initialize something ...

	greeterService := &greeterService{p}
	pb.RegisterGreeterServiceImp(p.Register, greeterService,
		transport.WithInterceptors(func(h interceptor.Handler) interceptor.Handler {
			fmt.Println("wrap greeterService methods")
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				info := ctx.Value(transport.ServiceInfoContextKey).(transport.ServiceInfo)
				fmt.Printf("before %s/%s\n", info.Service(), info.Method())
				fmt.Println(req)
				out, err := h(ctx, req)
				fmt.Printf("after %s/%s\n", info.Service(), info.Method())
				return out, err
			}
		}),
	)

	userService := &userService{p}
	pb.RegisterUserServiceImp(p.Register, userService,
		transport.WithInterceptors(func(h interceptor.Handler) interceptor.Handler {
			fmt.Println("wrap userService methods")
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				info := ctx.Value(transport.ServiceInfoContextKey).(transport.ServiceInfo)
				fmt.Printf("before %s/%s\n", info.Service(), info.Method())
				fmt.Println(req)
				out, err := h(ctx, req)
				fmt.Printf("after %s/%s\n", info.Service(), info.Method())
				return out, err
			}
		}),
	)

	return nil
}

func init() {
	servicehub.Register("erda.infra.example", &servicehub.Spec{
		Services:     pb.ServiceNames(),
		Types:        pb.Types(),
		Dependencies: []string{"service-register"},
		Description:  "",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
