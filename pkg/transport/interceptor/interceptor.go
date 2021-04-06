// Author: recallsong
// Email: songruiguo@qq.com

package interceptor

import "context"

// Handler defines the handler invoked by Interceptor.
type Handler func(ctx context.Context, req interface{}) (interface{}, error)

// Interceptor is HTTP/gRPC transport middleware.
type Interceptor func(Handler) Handler

// Chain returns a Interceptor that specifies the chained handler for endpoint.
func Chain(outer Interceptor, others ...Interceptor) Interceptor {
	if len(others) <= 0 {
		return outer
	}
	return func(next Handler) Handler {
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next)
		}
		return outer(next)
	}
}
