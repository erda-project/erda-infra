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

package httpendpoints

import (
	"context"
	"net/http"

	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/i18n"
	"github.com/gorilla/mux"
)

// Endpoint contains URL path and corresponding handler
type Endpoint struct {
	Path           string
	Method         string
	Handler        func(context.Context, *http.Request, map[string]string) (Responser, error)
	WriterHandler  func(context.Context, http.ResponseWriter, *http.Request, map[string]string) error
	ReverseHandler func(context.Context, *http.Request, map[string]string) error
}

// Interface .
type Interface interface {
	RegisterEndpoints(endpoints []Endpoint)
	Router() *mux.Router
}

// Responser is an interface for http response
type Responser interface {
	GetLocaledResp(locale i18n.LocaleResource) HTTPResponse
	GetStatus() int
	GetContent() interface{}
}
