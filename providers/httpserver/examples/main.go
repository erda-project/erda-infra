// Author: recallsong
// Email: songruiguo@qq.com

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

type define struct{}

func (d *define) Services() []string     { return []string{"hello"} }
func (d *define) Dependencies() []string { return []string{"http-server"} }
func (d *define) Description() string    { return "hello for example" }
func (d *define) Config() interface{}    { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type config struct {
	Message string `file:"message"`
}

type provider struct {
	Cfg *config
	Log logs.Logger
}

func (p *provider) Init(ctx servicehub.Context) error {
	// get httpserver.Router from service name "http-server"
	routes := ctx.Service("http-server",
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
	routes.Static("/hello/static", "/")
	routes.File("/hello/file", "/page.html")
	return nil
}

func init() {
	servicehub.RegisterProvider("hello", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
