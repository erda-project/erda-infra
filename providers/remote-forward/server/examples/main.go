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
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	forward "github.com/erda-project/erda-infra/providers/remote-forward"
	fserver "github.com/erda-project/erda-infra/providers/remote-forward/server"
)

type provider struct {
	FServer fserver.Interface `autowired:"remote-forward-server"`
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.FServer.AddHandshaker(func(req *forward.RequestHeader, resp *forward.ResponseHeader) error {
		resp.Values["response-data"] = "example"
		return nil
	})
	return nil
}

func init() {
	servicehub.Register("example", &servicehub.Spec{
		Creator: func() servicehub.Provider { return &provider{} },
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
