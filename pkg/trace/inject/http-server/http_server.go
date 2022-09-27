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
	"context"
	"net/http"
	_ "unsafe" //nolint

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	injectcontext "github.com/erda-project/erda-infra/pkg/trace/inject/context"
	"github.com/erda-project/erda-infra/pkg/trace/inject/hook"
)

type serverHandler struct {
	srv *http.Server
}

//go:linkname serveHTTP net/http.serverHandler.ServeHTTP
//go:noinline
func serveHTTP(s *serverHandler, rw http.ResponseWriter, req *http.Request)

//go:noinline
func originalServeHTTP(s *serverHandler, rw http.ResponseWriter, req *http.Request) {}

var tracedServerHandler = otelhttp.NewHandler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	injectcontext.SetContext(r.Context())
	defer injectcontext.ClearContext()
	s := getServerHandler(r.Context())
	originalServeHTTP(s, rw, r)
}), "", otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
	u := *r.URL
	u.RawQuery = ""
	u.ForceQuery = false
	return r.Method + " " + u.String()
}))

type _serverHandlerKey int8

const serverHandlerKey _serverHandlerKey = 0

func withServerHandler(ctx context.Context, s *serverHandler) context.Context {
	return context.WithValue(ctx, serverHandlerKey, s)
}

func getServerHandler(ctx context.Context) *serverHandler {
	return ctx.Value(serverHandlerKey).(*serverHandler)
}

//go:noinline
func wrappedHTTPHandler(s *serverHandler, rw http.ResponseWriter, req *http.Request) {
	req = req.WithContext(withServerHandler(req.Context(), s))
	tracedServerHandler.ServeHTTP(rw, req)
}

func init() {
	hook.Hook(serveHTTP, wrappedHTTPHandler, originalServeHTTP)
}
