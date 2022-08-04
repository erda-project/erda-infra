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
	"reflect"
	"sync"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/httpserver/interceptors"
	"github.com/erda-project/erda-infra/providers/httpserver/server"
)

// config .
type config struct {
	Addr        string `file:"addr" default:":8080" desc:"http address to listen"`
	PrintRoutes bool   `file:"print_routes" default:"true" desc:"print http routes"`
	AllowCORS   bool   `file:"allow_cors" default:"false" desc:"allow cors"`
	Reloadable  bool   `file:"reloadable" default:"false" desc:"routes reloadable"`

	Debug bool      `file:"debug" default:"false"`
	Log   LogConfig `file:"log"`
}

// LogConfig .
type LogConfig struct {
	MaxBodySizeBytes int `file:"max_body_size_bytes" default:"1024" desc:"max body size in bytes"`
}

type provider struct {
	Cfg *config
	Log logs.Logger

	server server.Server
	lock   sync.Mutex
	routes map[routeKey]*route
	err    error
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.server = server.New(p.Cfg.Reloadable, &dataBinder{}, &structValidator{validator: validator.New()})

	//p.server.Use(interceptors.Recover(p.Log).(func(echo.HandlerFunc) echo.HandlerFunc))
	p.server.Use(interceptors.SimpleRecord(p.getInterceptorOption()))
	p.server.Use(interceptors.CORS(p.Cfg.AllowCORS))
	p.server.Use(interceptors.InjectRequestID())
	p.server.Use(interceptors.DetailLog(p.getInterceptorOption()))
	p.server.Use(interceptors.BodyDump(p.getInterceptorOption(), p.Cfg.Log.MaxBodySizeBytes))
	p.server.Use(interceptors.PassThroughDebugFlag())
	p.server.Use(p.wrapContext())

	return nil
}

func (p *provider) wrapContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &context{Context: c}
			err := next(ctx)
			p.logFailure(c, err)
			return err
		}
	}
}

func (p *provider) getInterceptorOption() interceptors.Option {
	funcs := []interceptors.EnableFetchFunc{
		func(c echo.Context) bool {
			return p.Cfg.Debug
		},
	}
	return interceptors.NewOption(funcs, p.Log)
}

func (p *provider) logFailure(c server.Context, err error) {
	if err != nil || c.Response().Status/100 != 2 {
		p.Log.Errorf("(%s) err: %v, status: %d, url method: %s, path: %s, matcherPath: %s, ip: %s, header: %v",
			interceptors.GetRequestID(c), err, c.Response().Status, c.Request().Method, c.Request().URL.Path, c.Path(), c.RealIP(), c.Request().Header)
	}
}

// Start .
func (p *provider) Start() error {
	if p.err != nil {
		return p.err
	}
	if p.Cfg.PrintRoutes {
		if p.Cfg.Reloadable {
			p.lock.Lock()
		}
		p.printRoutes(p.routes)
		if p.Cfg.Reloadable {
			p.lock.Unlock()
		}
	}
	p.Log.Infof("starting http server at %s", p.Cfg.Addr)
	return p.server.Start(p.Cfg.Addr)
}

// Close .
func (p *provider) Close() error {
	if p.server == nil {
		return nil
	}
	return p.server.Close()
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Service() == "http-router-manager" || ctx.Type() == routerManagerType {
		return p.newRouterManager(true, ctx.Caller(), args...)
	} else if p.Cfg.Reloadable && (ctx.Service() != "http-router-tx" || ctx.Type() == routerType) {
		return &autoCommitRouter{
			tx: p.newRouterManager(false, ctx.Caller(), args...),
		}
	}
	return p.newRouterTx(true, ctx.Caller(), args...)
}

func (p *provider) newRouterManager(reset bool, group string, opts ...interface{}) RouterManager {
	return &routerManager{
		group: group,
		reset: reset,
		opts:  opts,
		p:     p,
	}
}

func (p *provider) newRouterTx(reset bool, group string, opts ...interface{}) RouterTx {
	interceptors := getInterceptors(opts)
	r := &router{
		tx:           p.server.NewRouter(),
		group:        group,
		interceptors: interceptors,
	}
	r.pathFormater = r.getPathFormater(opts)
	if p.Cfg.Reloadable {
		r.lock = &p.lock
		r.lock.Lock()
		r.routes = make(map[routeKey]*route)
		for key, route := range p.routes {
			if !reset || route.group != r.group {
				r.routes[key] = route
				if route.handler != nil {
					r.tx.Add(route.method, route.path, route.handler)
				}
			}
		}
		r.reportError = func(err error) {}
		r.updateRoutes = func(routes map[routeKey]*route) {
			p.routes = routes
			diff := make(map[routeKey]*route)
			for key, route := range p.routes {
				if route.group == r.group {
					diff[key] = route
				}
			}
			if p.Cfg.PrintRoutes {
				p.printRoutes(diff)
			}
		}
	} else {
		r.routes = p.routes
		r.reportError = func(err error) {
			p.err = err
		}
	}
	return r
}

var (
	routerType        = reflect.TypeOf((*Router)(nil)).Elem()
	routerTxType      = reflect.TypeOf((*RouterTx)(nil)).Elem()
	routerManagerType = reflect.TypeOf((*RouterManager)(nil)).Elem()
)

func init() {
	servicehub.Register("http-server", &servicehub.Spec{
		Services: []string{"http-server", "http-router", "http-router-manager", "http-router-tx"},
		Types: []reflect.Type{
			routerType,
			routerTxType,
			routerManagerType,
		},
		Description: "http server",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{
				routes: make(map[routeKey]*route),
			}
		},
	})
}
