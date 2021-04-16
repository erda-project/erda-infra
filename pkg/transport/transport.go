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

package transport

import (
	"context"

	transgrpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	transhttp "github.com/erda-project/erda-infra/pkg/transport/http"
	"github.com/erda-project/erda-infra/pkg/transport/interceptor"
)

// Register .
type Register interface {
	transhttp.Router
	transgrpc.ServiceRegistrar
}

// ServiceOption .
type ServiceOption func(*ServiceOptions)

// ServiceOptions .
type ServiceOptions struct {
	HTTP []transhttp.HandleOption
	GRPC []transgrpc.HandleOption
}

// WithInterceptors .
func WithInterceptors(o ...interceptor.Interceptor) ServiceOption {
	return func(opts *ServiceOptions) {
		if len(o) <= 0 {
			return
		}
		inter := interceptor.Chain(o[0], o[1:]...)
		opts.HTTP = append(opts.HTTP, transhttp.WithInterceptor(inter))
		opts.GRPC = append(opts.GRPC, transgrpc.WithInterceptor(inter))
	}
}

// WithHTTPOption .
func WithHTTPOptions(o ...transhttp.HandleOption) ServiceOption {
	return func(opts *ServiceOptions) {
		opts.HTTP = append(opts.HTTP, o...)
	}
}

// DefaultServiceOptions .
func DefaultServiceOptions() *ServiceOptions {
	return &ServiceOptions{}
}

// ServiceInfo .
type ServiceInfo interface {
	Service() string
	Method() string
	Instance() interface{}
}

// NewServiceInfo .
func NewServiceInfo(service string, method string, instance interface{}) ServiceInfo {
	return &serviceInfo{
		service:  service,
		method:   method,
		instance: instance,
	}
}

type serviceInfoContextKey int8

// ServiceInfoContextKey .
const ServiceInfoContextKey = serviceInfoContextKey(0)

// WithRequest .
func WithServiceInfo(ctx context.Context, info ServiceInfo) context.Context {
	return context.WithValue(ctx, ServiceInfoContextKey, info)
}

// ContextRequest .
func ContextServiceInfo(ctx context.Context) ServiceInfo {
	info, _ := ctx.Value(ServiceInfoContextKey).(ServiceInfo)
	return info
}

type serviceInfo struct {
	service  string
	method   string
	instance interface{}
}

func (si *serviceInfo) Service() string       { return si.service }
func (si *serviceInfo) Method() string        { return si.method }
func (si *serviceInfo) Instance() interface{} { return si.instance }
