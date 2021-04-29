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

package expvar

import (
	"expvar"
	"fmt"
	"net/http"
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/httpserver"
)

type config struct {
	Publish []string `file:"publish"`
}

// +provider
type provider struct {
	Cfg    *config
	Router httpserver.Router `autowired:"http-server@admin"`
}

// Run this is optional
func (p *provider) Init(ctx servicehub.Context) error {
	for _, item := range p.Cfg.Publish {
		v, ok := defaultVars[item]
		if !ok {
			return fmt.Errorf("var %q not exit", item)
		}
		expvar.Publish(item, v)
	}
	p.Router.Add(http.MethodGet, "/debug/vars", expvar.Handler())
	return nil
}

var defaultVars = map[string]expvar.Var{
	"envs": expvar.Func(func() interface{} {
		return os.Environ()
	}),
}

func init() {
	servicehub.Register("expvar", &servicehub.Spec{
		Services:    []string{"expvar"},
		Description: "expvar",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
