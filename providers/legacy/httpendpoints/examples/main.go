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

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints"
	_ "github.com/erda-project/erda-infra/providers/legacy/httpendpoints"
	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/errorresp"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Services() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"http-endpoints"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return an instance representing the configuration
func (d *define) Config() interface{} { return &config{} }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type config struct{}

type provider struct {
	C *config     // auto inject this field
	L logs.Logger // auto inject this field
}

func (p *provider) Init(ctx servicehub.Context) error {
	// register some apis
	server := ctx.Service("http-endpoints").(httpendpoints.Interface)
	server.RegisterEndpoints([]httpendpoints.Endpoint{
		httpendpoints.Endpoint{
			Path:    "/hello",
			Method:  http.MethodGet,
			Handler: p.Hello,
		},
		httpendpoints.Endpoint{
			Path:    "/error",
			Method:  http.MethodGet,
			Handler: p.Error,
		},
	})
	return nil
}

func (p *provider) Hello(ctx context.Context, r *http.Request, vars map[string]string) (
	httpendpoints.Responser, error) {
	return httpendpoints.OkResp(map[string]interface{}{
		"message": "ok",
	})
}

func (p *provider) Error(ctx context.Context, r *http.Request, vars map[string]string) (
	httpendpoints.Responser, error) {
	return errorresp.ErrResp(fmt.Errorf("example error"))
}

func init() {
	servicehub.RegisterProvider("example", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
