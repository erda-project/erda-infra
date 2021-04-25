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
