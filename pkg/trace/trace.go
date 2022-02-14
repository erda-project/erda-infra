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

package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	injectcontext "github.com/erda-project/erda-infra/pkg/trace/inject/context"
)

// Go .
func Go(ctx context.Context, fn func()) {
	GoWithContext(ctx, func(ctx context.Context) {
		fn()
	})
}

// GoWithContext .
func GoWithContext(ctx context.Context, fn func(ctx context.Context)) {
	if ctx == nil {
		ctx = context.Background()
	}
	if span := trace.SpanFromContext(ctx); !span.SpanContext().IsValid() {
		pctx := injectcontext.GetContext()
		if pctx != nil {
			if span := trace.SpanFromContext(pctx); span.SpanContext().IsValid() {
				ctx = trace.ContextWithSpan(ctx, span)
			}
		}
	}
	go injectcontext.RunWithContext(ctx, fn)
}

// ContextWithSpan .
func ContextWithSpan(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if span := trace.SpanFromContext(ctx); !span.SpanContext().IsValid() {
		pctx := injectcontext.GetContext()
		if pctx != nil {
			if span := trace.SpanFromContext(pctx); span.SpanContext().IsValid() {
				ctx = trace.ContextWithSpan(ctx, span)
			}
		}
	}
	return ctx
}

// RunWithTrace .
func RunWithTrace(ctx context.Context, fn func(ctx context.Context)) {
	injectcontext.RunWithContext(ctx, fn)
}
