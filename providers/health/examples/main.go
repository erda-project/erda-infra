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
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/health"
	_ "github.com/erda-project/erda-infra/providers/httpserver"
)

type provider struct {
}

func (p *provider) Init(ctx servicehub.Context) error {
	h := ctx.Service("health").(health.Interface)
	h.Register(p.HealthCheck)
	return nil
}

func (p *provider) HealthCheck(context.Context) error {
	return fmt.Errorf("error message")
}

func init() {
	servicehub.Register("examples", &servicehub.Spec{
		Services:     []string{"hello"},
		Dependencies: []string{"health"},
		Description:  "hello for example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
