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

package httpserver

import (
	"net/http"
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// config .
type config struct {
	Addr        string `file:"addr" default:":8080" desc:"http address to listen"`
	PrintRoutes bool   `file:"print_routes" default:"true" desc:"print http routes"`
	AllowCORS   bool   `file:"allow_cors" default:"false" desc:"allow cors"`
}

type provider struct {
	Cfg    *config
	Log    logs.Logger
	server *echo.Echo
	router *router
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.server = echo.New()
	p.server.HideBanner = true
	p.server.HidePort = true
	p.server.Binder = &dataBinder{}
	p.server.Validator = &structValidator{validator: validator.New()}
	if p.Cfg.AllowCORS {
		p.server.Use(middleware.CORS())
	}
	p.server.Use(func(fn echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx = &context{Context: ctx}
			err := fn(ctx)
			if err != nil {
				p.Log.Error(err)
				return err
			}
			return nil
		}
	})
	return nil
}

// Start .
func (p *provider) Start() error {
	if p.Cfg.PrintRoutes /*|| p.Cfg.IndexShowRoutes*/ {
		p.router.Normalize()
	}
	if p.Cfg.PrintRoutes {
		for _, route := range p.router.routes {
			if !route.hide {
				p.Log.Infof("%s --> %s", p.Cfg.Addr, route.String())
			}
		}
	}
	p.Log.Infof("starting http server at %s", p.Cfg.Addr)
	err := p.server.Start(p.Cfg.Addr)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Close .
func (p *provider) Close() error {
	if p.server == nil || p.server.Server == nil {
		return nil
	}
	err := p.server.Server.Close()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	interceptors := getInterceptors(args)
	r := &router{
		p:            p,
		routeMap:     p.router.routeMap,
		group:        ctx.Caller(),
		interceptors: interceptors,
	}
	r.pathFormater = r.getPathFormater(args)
	return Router(r)
}

func init() {
	servicehub.Register("http-server", &servicehub.Spec{
		Services:    []string{"http-server", "http-routes", "http-router"},
		Types:       []reflect.Type{reflect.TypeOf((*Router)(nil)).Elem()},
		Description: "http server",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			p := &provider{
				router: &router{
					routeMap: make(map[routeKey]*route),
				},
			}
			p.router.p = p
			return p
		},
	})
}
