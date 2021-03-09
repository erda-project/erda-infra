// Author: recallsong
// Email: songruiguo@qq.com

package httpserver

import (
	"net/http"

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

type providerDefine struct{}

func (d *providerDefine) Service() []string {
	return []string{"http-server", "api-server"}
}

func (d *providerDefine) Summary() string {
	return "http server"
}

func (d *providerDefine) Description() string {
	return d.Summary()
}

func (d *providerDefine) Creator() servicehub.Creator {
	return newProvider
}

func (d *providerDefine) Config() interface{} {
	return &config{}
}

type provider struct {
	Cfg    *config
	Logger logs.Logger
	server *echo.Echo
	router *router
}

func newProvider() servicehub.Provider {
	p := &provider{
		router: &router{
			routeMap: make(map[routeKey]*route),
		},
	}
	p.router.p = p
	return p
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
				p.Logger.Error(err)
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
				p.Logger.Infof("%s --> %s", p.Cfg.Addr, route.String())
			}
		}
	}
	p.Logger.Infof("starting http server at %s", p.Cfg.Addr)
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
func (p *provider) Provide(name string, args ...interface{}) interface{} {
	interceptors := getInterceptors(args)
	return Router(&router{
		p:            p,
		routeMap:     p.router.routeMap,
		group:        name,
		interceptors: interceptors,
	})
}

func init() {
	servicehub.RegisterProvider("http-server", &providerDefine{})
}
