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

	"github.com/xormplus/xorm"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/mysqlxorm"
)

type provider struct {
	DB    *xorm.Engine        // autowired
	MySQL mysqlxorm.Interface // autowired
}

func (p *provider) Init(ctx servicehub.Context) error {
	fmt.Println(p.DB)
	fmt.Println(p.MySQL)
	// do something
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	r, err := p.MySQL.DB().QueryString("show tables")
	if err != nil {
		panic(err)
	}
	for i, m := range r {
		for k, v := range m {
			fmt.Println(i, k, v)
		}
	}
	return nil
}

func init() {
	servicehub.Register("example", &servicehub.Spec{
		Services:     []string{"example"},
		Dependencies: []string{"mysql-xorm"},
		Description:  "example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
