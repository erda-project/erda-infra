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
	"context"
	"net/http"
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	_ "github.com/erda-project/erda-infra/providers/health"
	"github.com/erda-project/erda-infra/providers/httpserver"
)

type provider struct {
	Log logs.Logger
	rm  httpserver.RouterManager
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.rm = ctx.Service("http-router-manager",
		func(handler func(ctx httpserver.Context) error) func(ctx httpserver.Context) error {
			return func(ctx httpserver.Context) error {
				p.Log.Info("intercept request", ctx.Request().URL.String())
				return handler(ctx)
			}
		},
	).(httpserver.RouterManager)

	routes := p.rm.NewRouter()
	routes.GET("/hello",
		func(resp http.ResponseWriter, req *http.Request) {
			resp.Write([]byte("hello"))
		},
	)

	return routes.Commit()
}

func (p *provider) Run(ctx context.Context) error {
	select {
	case <-time.After(5 * time.Second):
		p.Log.Infof("routes reload")
		r := p.rm.NewRouter(func(handler func(ctx httpserver.Context) error) func(ctx httpserver.Context) error {
			return func(ctx httpserver.Context) error {
				p.Log.Info("new common intercept request", ctx.Request().URL.String())
				return handler(ctx)
			}
		})
		r.GET("/hello3",
			func(resp http.ResponseWriter, req *http.Request) {
				resp.Write([]byte("hello3"))
			},
			func(handler func(ctx httpserver.Context) error) func(ctx httpserver.Context) error {
				return func(ctx httpserver.Context) error {
					p.Log.Info("new api intercept request", ctx.Request().URL.String())
					return handler(ctx)
				}
			},
		)
		return r.Commit()
	case <-ctx.Done():
	}
	return nil
}

type provider2 struct {
	Log logs.Logger
}

func (p *provider2) Init(ctx servicehub.Context) error {
	rm := ctx.Service("http-router-manager",
		func(handler func(ctx httpserver.Context) error) func(ctx httpserver.Context) error {
			return func(ctx httpserver.Context) error {
				p.Log.Info("intercept request", ctx.Request().URL.String())
				return handler(ctx)
			}
		},
	).(httpserver.RouterManager)

	routes := rm.NewRouter()
	routes.GET("/hello2",
		func(resp http.ResponseWriter, req *http.Request) {
			resp.Write([]byte("hello2"))
		},
	)
	return routes.Commit()
}

func init() {
	servicehub.Register("hello", &servicehub.Spec{
		Dependencies: []string{"http-server"},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
	servicehub.Register("hello2", &servicehub.Spec{
		Dependencies: []string{"http-server"},
		Creator: func() servicehub.Provider {
			return &provider2{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
