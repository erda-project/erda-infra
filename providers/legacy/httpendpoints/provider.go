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

package httpendpoints

import (
	"net/http"
	"reflect"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/i18n"
	"github.com/gorilla/mux"
)

type define struct{}

func (d *define) Services() []string { return []string{"http-endpoints"} }
func (d *define) Types() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Interface)(nil)).Elem()}
}
func (d *define) Dependencies() []string { return []string{"i18n"} }
func (d *define) Description() string    { return "http endpoints" }
func (d *define) Config() interface{} {
	return &config{}
}
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{
			router: mux.NewRouter(),
		}
	}
}

// config .
type config struct {
	Addr string `file:"addr" default:":8090" desc:"http address to listen"`
}

var _ Interface = (*provider)(nil)

type provider struct {
	C      *config
	L      logs.Logger
	router *mux.Router
	srv    *http.Server
	t      i18n.Translator
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	i := ctx.Service("i18n").(i18n.I18n)
	p.t = i.Translator("httpendpoints")
	p.srv = &http.Server{
		Addr:              p.C.Addr,
		Handler:           p.router,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
	}
	return nil
}

// Start .
func (p *provider) Start() error {
	p.L.Infof("starting endpoints at %s", p.C.Addr)
	return p.srv.ListenAndServe()
}

func (p *provider) Router() *mux.Router { return p.router }

// Close .
func (p *provider) Close() error {
	return p.srv.Close()
}

func init() {
	servicehub.RegisterProvider("http-endpoints", &define{})
}
