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

package traceinject

import (
	"net/http"
	"time"
	_ "unsafe" //nolint

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"

	injectcontext "github.com/erda-project/erda-infra/pkg/trace/inject/context"
	"github.com/erda-project/erda-infra/pkg/trace/inject/hook"
)

//go:linkname send net/http.send
//go:noinline
func send(ireq *http.Request, rt http.RoundTripper, deadline time.Time) (resp *http.Response, didTimeout func() bool, err error)

//go:noinline
func originalSend(ireq *http.Request, rt http.RoundTripper, deadline time.Time) (resp *http.Response, didTimeout func() bool, err error) {
	return send(ireq, rt, deadline)
}

//go:noinline
func tracedSend(ireq *http.Request, rt http.RoundTripper, deadline time.Time) (resp *http.Response, didTimeout func() bool, err error) {
	rt = &wrappedRoundTripper{
		RoundTripper: otelhttp.NewTransport(rt),
	}
	return originalSend(ireq, rt, deadline)
}

type wrappedRoundTripper struct {
	http.RoundTripper
}

//go:noinline
func (t *wrappedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = contextWithSpan(req)
	return t.RoundTripper.RoundTrip(req)
}

//go:noinline
func contextWithSpan(req *http.Request) *http.Request {
	ctx := req.Context()
	if span := trace.SpanFromContext(ctx); !span.SpanContext().IsValid() {
		pctx := injectcontext.GetContext()
		if pctx != nil {
			if span := trace.SpanFromContext(pctx); span.SpanContext().IsValid() {
				ctx = trace.ContextWithSpan(ctx, span)
				req = req.WithContext(ctx)
			}
		}
	}
	return req
}

func init() {
	hook.Hook(send, tracedSend, originalSend)
}
