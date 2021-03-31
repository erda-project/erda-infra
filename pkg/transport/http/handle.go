// Author: recallsong
// Email: songruiguo@qq.com

package http

import (
	"net/http"

	"github.com/erda-project/erda-infra/pkg/transport/interceptor"
)

const (
	// SupportPackageIsVersion1 These constants should not be referenced from any other code.
	SupportPackageIsVersion1 = true
)

// DecodeRequestFunc is decode request func.
type DecodeRequestFunc func(*http.Request, interface{}) error

// EncodeResponseFunc is encode response func.
type EncodeResponseFunc func(http.ResponseWriter, *http.Request, interface{}) error

// EncodeErrorFunc is encode error func.
type EncodeErrorFunc func(http.ResponseWriter, *http.Request, error)

// HandleOption is handle option.
type HandleOption func(*HandleOptions)

// HandleOptions is handle options.
type HandleOptions struct {
	Decode      DecodeRequestFunc
	Encode      EncodeResponseFunc
	Error       EncodeErrorFunc
	Interceptor interceptor.Interceptor
}

// HandlerFunc .
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Router .
type Router interface {
	Add(method, path string, handler HandlerFunc)
}

// RouterFunc .
type RouterFunc func(method, path string, handler HandlerFunc)

// Add .
func (fn RouterFunc) Add(method, path string, handler HandlerFunc) {
	fn(method, path, handler)
}
