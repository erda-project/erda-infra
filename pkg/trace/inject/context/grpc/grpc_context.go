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

package grpccontext

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	injectcontext "github.com/erda-project/erda-infra/pkg/trace/inject/context"
)

// UnaryServerInterceptor .
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	inter := otelgrpc.UnaryServerInterceptor()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return inter(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
			injectcontext.SetContext(ctx)
			defer injectcontext.ClearContext()
			return handler(ctx, req)
		})
	}
}

// StreamServerInterceptor .
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	inter := otelgrpc.StreamServerInterceptor()
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return inter(srv, ss, info, func(srv interface{}, stream grpc.ServerStream) error {
			injectcontext.SetContext(stream.Context())
			defer injectcontext.ClearContext()
			return handler(srv, stream)
		})
	}
}

// UnaryClientInterceptor .
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	inter := otelgrpc.UnaryClientInterceptor()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = injectcontext.ContextWithSpan(ctx)
		return inter(ctx, method, req, reply, cc, invoker, opts...)
	}
}

// StreamClientInterceptor .
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	inter := otelgrpc.StreamClientInterceptor()
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = injectcontext.ContextWithSpan(ctx)
		return inter(ctx, desc, cc, method, streamer, opts...)
	}
}
