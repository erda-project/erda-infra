package transport

import (
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

type serviceInfo struct {
	service  string
	method   string
	instance interface{}
}

func (si *serviceInfo) Service() string       { return si.service }
func (si *serviceInfo) Method() string        { return si.method }
func (si *serviceInfo) Instance() interface{} { return si.instance }
