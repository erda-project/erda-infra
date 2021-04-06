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

// WithInterceptor .
func WithInterceptor(o interceptor.Interceptor) HandleOption {
	return func(opts *HandleOptions) {
		opts.Interceptor = o
	}
}

// WithDecoder .
func WithDecoder(o DecodeRequestFunc) HandleOption {
	return func(opts *HandleOptions) {
		opts.Decode = o
	}
}

// WithEncoder .
func WithEncoder(o EncodeResponseFunc) HandleOption {
	return func(opts *HandleOptions) {
		opts.Encode = o
	}
}

// WithErrorEncoder .
func WithErrorEncoder(o EncodeErrorFunc) HandleOption {
	return func(opts *HandleOptions) {
		opts.Error = o
	}
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

type requestContextKey int8

// RequestContextKey .
const RequestContextKey = requestContextKey(0)
