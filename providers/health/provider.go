// Author: recallsong
// Email: songruiguo@qq.com

package health

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/httpserver"
)

// Checker .
type Checker func() error

// Interface .
type Interface interface {
	Register(Checker)
}

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
	return func() servicehub.Provider {
		return &provider{
			checkers: make(map[string][]Checker),
		}
	}
}

type provider struct {
	C        *config
	names    []string
	checkers map[string][]Checker
	body     []byte
}

func (p *provider) Init(ctx servicehub.Context) error {
	routes := ctx.Service("http-server").(httpserver.Router)
	for _, path := range p.C.Path {
		routes.GET(path, p.handler)
	}
	p.body = []byte(p.C.Body)
	return nil
}

func (p *provider) handler(resp http.ResponseWriter, req *http.Request) error {
	for _, key := range p.names {
		for _, checker := range p.checkers[key] {
			err := checker()
			if err != nil {
				return fmt.Errorf("%s is unhealthy: %s", key, err)
			}
		}
	}
	resp.Header().Set("Content-Type", p.C.ContentType)
	resp.WriteHeader(p.C.Status)
	resp.Write(p.body)
	return nil
}

// Provide .
func (p *provider) Provide(name string, args ...interface{}) interface{} {
	return &service{
		name: name,
		p:    p,
	}
}

type service struct {
	name string
	p    *provider
}

func (s *service) Register(c Checker) {
	list, ok := s.p.checkers[s.name]
	if !ok {
		s.p.names = append(s.p.names, s.name)
		sort.Strings(s.p.names)
	}
	s.p.checkers[s.name] = append(list, c)
}

func init() {
	servicehub.RegisterProvider("health", &define{})
}
