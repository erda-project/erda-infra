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
	_ "unsafe" //nolint

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"

	injectcontext "github.com/erda-project/erda-infra/pkg/trace/inject/context"
	"github.com/erda-project/erda-infra/pkg/trace/inject/hook"
)

//go:linkname RoundTrip net/http.(*Transport).RoundTrip
//go:noinline
// RoundTrip .
func RoundTrip(t *http.Transport, req *http.Request) (*http.Response, error)

//go:noinline
func originalRoundTrip(t *http.Transport, req *http.Request) (*http.Response, error) {
	return RoundTrip(t, req)
}

type wrappedTransport struct {
	t *http.Transport
}

//go:noinline
func (t *wrappedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return originalRoundTrip(t.t, req)
}

//go:noinline
func tracedRoundTrip(t *http.Transport, req *http.Request) (*http.Response, error) {
	req = contextWithSpan(req)
	return otelhttp.NewTransport(&wrappedTransport{t: t}).RoundTrip(req)
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
	hook.Hook(RoundTrip, tracedRoundTrip, originalRoundTrip)
}
