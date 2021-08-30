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

package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/httpserver"
)

type config struct {
	MetricsPath string `file:"metrics_path" default:"/metrics"`
}

// provider .
type provider struct {
	server *http.Server
	Cfg    *config
	Router httpserver.Router `autowire:"http-server@admin"`
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.Router.GET(p.Cfg.MetricsPath, promhttp.Handler())
	return nil
}

func init() {
	servicehub.Register("prometheus", &servicehub.Spec{
		Services:     []string{"prometheus"},
		Dependencies: []string{"http-server"},
		Description:  "bind prometheus endpoint to http-server",
		ConfigFunc:   func() interface{} { return &config{} },
		Creator:      func() servicehub.Provider { return &provider{} },
	})
}
