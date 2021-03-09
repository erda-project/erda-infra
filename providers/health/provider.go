// Author: recallsong
// Email: songruiguo@qq.com

package health

import (
	"net/http"

	"github.com/erda-project/erda-infra/providers/httpserver"
	"github.com/erda-project/erda-infra/base/servicehub"
)

type config struct {
	Path        []string `file:"path" default:"/health" desc:"http path"`
	Status      int      `file:"status" default:"200" desc:"http response status"`
	Body        string   `file:"body" default:"{\"success\":true,\"data\":\"ok\"}" desc:"http response body"`
	ContentType string   `file:"content_type" default:"application/json" desc:"http response Content-Type"`
}

type define struct{}

func (d *define) Service() []string      { return []string{"health", "health-check"} }
func (d *define) Dependencies() []string { return []string{"http-server"} }
func (d *define) Description() string    { return "http health check" }
func (d *define) Config() interface{}    { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider { return &provider{} }
}

type provider struct {
	C    *config
	body []byte
}

func (p *provider) Init(ctx servicehub.Context) error {
	routes := ctx.Service("http-server").(httpserver.Router)
	for _, path := range p.C.Path {
		routes.GET(path, p.handler)
	}
	p.body = []byte(p.C.Body)
	return nil
}

func (p *provider) handler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", p.C.ContentType)
	resp.WriteHeader(p.C.Status)
	resp.Write(p.body)
}

func init() {
	servicehub.RegisterProvider("health", &define{})
}
