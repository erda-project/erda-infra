// Author: recallsong
// Email: songruiguo@qq.com

package pprof

import (
	"net/http"
	"net/http/pprof"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/httpserver"
)

type define struct{}

func (d *define) Services() []string     { return []string{"pprof"} }
func (d *define) Dependencies() []string { return []string{"http-server"} }
func (d *define) Summary() string        { return "start pprof http server" }
func (d *define) Description() string    { return d.Summary() }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider { return &provider{} }
}

// provider .
type provider struct {
	server *http.Server
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	routes := ctx.Service("http-server@admin").(httpserver.Router)
	routes.Any("/debug/pprof/**", mux)
	return nil
}

func init() {
	servicehub.RegisterProvider("pprof", &define{})
}
