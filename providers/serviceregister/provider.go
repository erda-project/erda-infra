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

package register

import (
	"fmt"
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"
	transgrpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	transhttp "github.com/erda-project/erda-infra/pkg/transport/http"
	"github.com/erda-project/erda-infra/providers/grpcserver"
	"github.com/erda-project/erda-infra/providers/httpserver"
	"google.golang.org/grpc"
)

// Interface .
type Interface = transport.Register

type provider struct {
	Log    logs.Logger
	router httpserver.Router
	grpc   grpcserver.Interface
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.router, _ = ctx.Service("http-server").(httpserver.Router)
	p.grpc, _ = ctx.Service("grpc-server").(grpcserver.Interface)
	if p.router == nil && p.grpc == nil {
		return fmt.Errorf("not found http-server or grpc-server")
	}
	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return &service{
		name:   ctx.Caller(),
		router: p.router,
		grpc:   p.grpc,
	}
}

var _ Interface = (*service)(nil)

type service struct {
	name   string
	router httpserver.Router
	grpc   grpcserver.Interface
}

func (s *service) Add(method, path string, handler transhttp.HandlerFunc) {
	if s.router != nil {
		s.router.Add(method, path, handler, httpserver.WithPathFormat(httpserver.PathFormatGoogleAPIs))
	}
}

func (s *service) RegisterService(sd *grpc.ServiceDesc, impl interface{}) {
	if s.grpc != nil {
		s.grpc.RegisterService(sd, impl)
	}
}

func init() {
	servicehub.Register("service-register", &servicehub.Spec{
		Services: []string{"service-register"},
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
			reflect.TypeOf((*transgrpc.ServiceRegistrar)(nil)).Elem(),
			reflect.TypeOf((*transhttp.Router)(nil)).Elem(),
		},
		OptionalDependencies: []string{"grpc-server", "http-server"},
		Description:          "provide grpc and http server",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
