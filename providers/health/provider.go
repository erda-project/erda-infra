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

package health

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"sort"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/httpserver"
)

// Checker .
type Checker func(context.Context) error

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

type provider struct {
	Cfg          *config
	Router       httpserver.Router `autowire:"http-server"`
	names        []string
	checkers     map[string][]Checker
	healthBody   []byte
	unhealthBody []byte
}

func (p *provider) Init(ctx servicehub.Context) error {
	for _, path := range p.Cfg.Path {
		p.Router.GET(path, p.handler)
	}
	p.healthBody = []byte(p.Cfg.HealthBody)
	p.unhealthBody = []byte(p.Cfg.UnhealthBody)
	return nil
}

func (p *provider) handler(resp http.ResponseWriter, req *http.Request) error {
	status := make(map[string]interface{})
	health := true
	for _, key := range p.names {
		var errors []interface{}
		for _, checker := range p.checkers[key] {
			err := checker(context.Background())
			if err != nil {
				errors = append(errors, err.Error())
				health = false
				if p.Cfg.AbortOnError {
					break
				}
			}
		}
		status[key] = errors
	}
	resp.Header().Set("Content-Type", p.Cfg.ContentType)
	var body []byte
	if health {
		resp.WriteHeader(p.Cfg.HealthStatus)
		body = p.healthBody
	} else {
		resp.WriteHeader(p.Cfg.UnhealthStatus)
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
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return &service{
		name: ctx.Caller(),
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
	servicehub.Register("health", &servicehub.Spec{
		Services:     []string{"health", "health-checker"},
		Types:        []reflect.Type{reflect.TypeOf((*Interface)(nil)).Elem()},
		Dependencies: []string{"http-server"},
		Description:  "http health check",
		ConfigFunc:   func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{
				checkers: make(map[string][]Checker),
			}
		},
	})
}
