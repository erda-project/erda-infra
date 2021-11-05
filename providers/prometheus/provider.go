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
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/httpserver"
)

type config struct {
	MetricsPath string `file:"metrics_path" default:"/metrics"`
	HTTPRouter  string `file:"http_router" default:"http-server@admin"`
}

// provider .
type provider struct {
	server *http.Server
	Cfg    *config
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	svc := ctx.Service(p.Cfg.HTTPRouter)
	if svc == nil {
		return errors.New("unable to find http router: "+ p.Cfg.HTTPRouter)
	}
	router, ok := svc.(httpserver.Router)
	if !ok {
		return fmt.Errorf("invalid type %T, which must be httpserver.Router", svc)
	}
	router.GET(p.Cfg.MetricsPath, promhttp.Handler())
	return nil
}

func init() {
	servicehub.Register("prometheus", &servicehub.Spec{
		Services:     []string{"prometheus"},
		Description:  "bind prometheus endpoint to http-server",
		Dependencies: []string{"http-server"},
		ConfigFunc:   func() interface{} { return &config{} },
		Creator:      func() servicehub.Provider { return &provider{} },
	})
}
