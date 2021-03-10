// Author: recallsong
// Email: songruiguo@qq.com

package health

import (
	"encoding/json"
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
	Path           []string `file:"path" default:"/health" desc:"http path"`
	HealthStatus   int      `file:"health_status" default:"200" desc:"http response status if health"`
	UnhealthStatus int      `file:"unhealth_status" default:"503" desc:"http response status if unhealth"`
	HealthBody     string   `file:"health_body" desc:"http response body if health"`
	UnhealthBody   string   `file:"unhealth_body" desc:"http response body if unhealth"`
	ContentType    string   `file:"content_type" default:"application/json" desc:"http response Content-Type"`
	AbortOnError   bool     `file:"abort_on_error"`
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
	C            *config
	names        []string
	checkers     map[string][]Checker
	healthBody   []byte
	unhealthBody []byte
}

func (p *provider) Init(ctx servicehub.Context) error {
	routes := ctx.Service("http-server").(httpserver.Router)
	for _, path := range p.C.Path {
		routes.GET(path, p.handler)
	}
	p.healthBody = []byte(p.C.HealthBody)
	p.unhealthBody = []byte(p.C.UnhealthBody)
	return nil
}

func (p *provider) handler(resp http.ResponseWriter, req *http.Request) error {
	status := make(map[string]interface{})
	health := true
	for _, key := range p.names {
		var errors []interface{}
		for _, checker := range p.checkers[key] {
			err := checker()
			if err != nil {
				errors = append(errors, err.Error())
				health = false
				if p.C.AbortOnError {
					break
				}
			}
		}
		status[key] = errors
	}
	resp.Header().Set("Content-Type", p.C.ContentType)
	var body []byte
	if health {
		resp.WriteHeader(p.C.HealthStatus)
		body = p.healthBody
	} else {
		resp.WriteHeader(p.C.UnhealthStatus)
		body = p.unhealthBody
	}
	if len(body) > 0 {
		resp.Write(body)
	} else {
		byts, _ := json.Marshal(map[string]interface{}{
			"health":   health,
			"checkers": status,
		})
		resp.Write(byts)
	}
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
