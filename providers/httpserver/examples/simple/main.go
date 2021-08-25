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

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/httpserver"
)

type config struct {
	Message string `file:"message"`
}

type provider struct {
	Cfg *config
	Log logs.Logger
}

func (p *provider) Init(ctx servicehub.Context) error {
	// get httpserver.Router from service name "http-router"
	routes := ctx.Service("http-router",
		// this is interceptor for this provider
		func(handler func(ctx httpserver.Context) error) func(ctx httpserver.Context) error {
			return func(ctx httpserver.Context) error {
				fmt.Println("intercept request", ctx.Request().URL.String())
				return handler(ctx)
			}
		},
	).(httpserver.Router)

	// request parameters http.ResponseWriter, *http.Request
	routes.GET("/hello",
		func(resp http.ResponseWriter, req *http.Request) {
			resp.Write([]byte(p.Cfg.Message))
		},
		httpserver.WithDescription("this is hello provider"),
		httpserver.WithInterceptor(
			func(handler func(ctx httpserver.Context) error) func(ctx httpserver.Context) error {
				return func(ctx httpserver.Context) error {
					return handler(ctx)
				}
			},
		),
	)

	// request parameter is struct pointer, response is: status int, data interface{}, err error
	routes.POST("/hello/simple", func(body *struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}) (status int, data interface{}, err error) {
		return http.StatusCreated, body, nil
	})

	// request parameter is struct, and validate message field
	routes.POST("/hello/struct/:name", func(resp http.ResponseWriter, req *http.Request,
		body struct {
			Name    string `param:"name"`
			Message string `json:"message" form:"message" query:"message" validate:"required"`
		},
	) {
		resp.Write([]byte(fmt.Sprint(body)))
	})

	// request parameter is struct
	routes.POST("/hello/struct/ptr/:name", func(resp http.ResponseWriter, req *http.Request,
		body *struct {
			Name    string `param:"name"`
			Message string `json:"message" form:"message" query:"message" validate:"required"`
		},
	) {
		resp.Write([]byte(fmt.Sprint(body)))
	})

	// request parameters: http.ResponseWriter, *http.Request, []byte, and []byte is request Body
	routes.Any("/hello/bytes", func(resp http.ResponseWriter, req *http.Request, byts []byte) {
		resp.Write(byts)
	})

	// request parameters: http.ResponseWriter, *http.Request, int
	routes.Any("/hello/int", func(resp http.ResponseWriter, req *http.Request, body int) {
		resp.Write([]byte(fmt.Sprint(body)))
	})
	routes.Any("/hello/int/ptr", func(resp http.ResponseWriter, req *http.Request, body *int) {
		resp.Write([]byte(fmt.Sprint(*body)))
	})

	// request parameters: http.ResponseWriter, *http.Request, map[string]interface{}
	routes.Any("/hello/map", func(resp http.ResponseWriter, req *http.Request, body map[string]interface{}) {
		resp.Write([]byte(fmt.Sprint(body)))
	})
	routes.Any("/hello/map/ptr", func(resp http.ResponseWriter, req *http.Request, body ******map[string]interface{}) {
		resp.Write([]byte(fmt.Sprint(*body)))
	})

	// request parameters: http.ResponseWriter, *http.Request, []interface{}
	routes.Any("/hello/slice", func(resp http.ResponseWriter, req *http.Request, body []interface{}) {
		resp.Write([]byte(fmt.Sprint(body)))
	})

	// request parameters: httpserver.Context, string
	routes.POST("/hello/context", func(ctx httpserver.Context, body string) {
		ctx.ResponseWriter().Write([]byte(body))
	})

	// request parameters: status int, body io.Reader
	routes.GET("/hello/response/body", func(ctx httpserver.Context) (status int, body io.Reader) {
		return http.StatusOK, bytes.NewReader([]byte("hello"))
	})

	// handle static files
	routes.Static("/hello/static", "")
	routes.File("/hello/file", "index.html")
	return nil
}

func init() {
	servicehub.Register("hello", &servicehub.Spec{
		Services:     []string{"hello"},
		Dependencies: []string{"http-server"},
		Description:  "hello for example",
		ConfigFunc:   func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
