// Author: recallsong
// Email: songruiguo@qq.com

package grpc

import (
	"context"

	"github.com/erda-project/erda-infra/pkg/transport/interceptor"
	"google.golang.org/grpc"
)

// ServiceRegistrar wraps a single method that supports service registration. It
// enables users to pass concrete types other than grpc.Server to the service
// registration methods exported by the IDL generated code.
//
// Upward compatible for google.golang.org/grpc v1.28.0
// and ServiceRegistrar define in google.golang.org/grpc v1.32.0+
type ServiceRegistrar interface {
	// RegisterService registers a service and its implementation to the
	// concrete type implementing this interface.  It may not be called
	// once the server has started serving.
	// desc describes the service and its methods and handlers. impl is the
	// service implementation which is passed to the method handlers.
	RegisterService(desc *grpc.ServiceDesc, impl interface{})
}

// ClientConnInterface defines the functions clients need to perform unary and
// streaming RPCs.  It is implemented by *ClientConn, and is only intended to
// be referenced by generated code.
//
// Upward compatible for google.golang.org/grpc v1.26.0
// and ClientConnInterface define in google.golang.org/grpc v1.28.0+
type ClientConnInterface interface {
	// Invoke performs a unary RPC and returns after the response is received
	// into reply.
	Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error
	// NewStream begins a streaming RPC.
	NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
}

// HandleOption is handle option.
type HandleOption func(*HandleOptions)

// HandleOptions is handle options.
type HandleOptions struct {
	Interceptor interceptor.Interceptor
}

// DefaultHandleOptions .
func DefaultHandleOptions() *HandleOptions {
	return &HandleOptions{}
}
