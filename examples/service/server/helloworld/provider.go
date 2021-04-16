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

package example

import (
	"context"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/examples/service/protocol/pb"
	"github.com/erda-project/erda-infra/pkg/transport"
	"github.com/erda-project/erda-infra/pkg/transport/interceptor"
)

type config struct {
}

// +provider
type provider struct {
	Cfg            *config
	Log            logs.Logger
	Register       transport.Register
	greeterService *greeterService
	userService    *userService
}

func (p *provider) Init(ctx servicehub.Context) error {
	// TODO initialize something ...

	if p.Register != nil {
		p.greeterService = &greeterService{p}
		pb.RegisterGreeterServiceImp(p.Register, p.greeterService,
			transport.WithInterceptors(func(h interceptor.Handler) interceptor.Handler {
				p.Log.Info("wrap greeterService methods")
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					info := transport.ContextServiceInfo(ctx)
					p.Log.Infof("before %s/%s\n", info.Service(), info.Method())
					p.Log.Info(req)
					out, err := h(ctx, req)
					p.Log.Infof("after %s/%s\n", info.Service(), info.Method())
					return out, err
				}
			}),
		)

		p.userService = &userService{p}
		pb.RegisterUserServiceImp(p.Register, p.userService,
			transport.WithInterceptors(func(h interceptor.Handler) interceptor.Handler {
				p.Log.Info("wrap userService methods")
				return func(ctx context.Context, req interface{}) (interface{}, error) {
					info := transport.ContextServiceInfo(ctx)
					p.Log.Infof("before %s/%s\n", info.Service(), info.Method())
					p.Log.Info(req)
					out, err := h(ctx, req)
					p.Log.Infof("after %s/%s\n", info.Service(), info.Method())
					return out, err
				}
			}),
		)
	}
	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	switch {
	case ctx.Service() == "erda.infra.example.GreeterService" || ctx.Type() == pb.GreeterServiceServerType() || ctx.Type() == pb.GreeterServiceHandlerType():
		return p.greeterService
	case ctx.Service() == "erda.infra.example.UserService" || ctx.Type() == pb.UserServiceServerType() || ctx.Type() == pb.UserServiceHandlerType():
		return p.userService
	}
	return p
}

func init() {
	servicehub.Register("erda.infra.example", &servicehub.Spec{
		Services:             pb.ServiceNames(),
		Types:                pb.Types(),
		OptionalDependencies: []string{"service-register"},
		Description:          "",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
