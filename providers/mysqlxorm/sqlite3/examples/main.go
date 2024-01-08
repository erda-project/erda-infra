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
	"github.com/erda-project/erda-infra/providers/mysqlxorm"
	_ "github.com/erda-project/erda-infra/providers/mysqlxorm/sqlite3"
)

type provider struct {
	Sqlite3 mysqlxorm.Interface
}

func (p *provider) Init(ctx servicehub.Context) error {
	fmt.Println(p.Sqlite3)

	return nil
}

func (p *provider) Run(ctx context.Context) error {
	err := p.Sqlite3.DB().Ping()
	if err != nil {
		fmt.Printf("connect sqlite3 error : %s \n", err)
	}
	return nil
}

func init() {
	servicehub.Register("example", &servicehub.Spec{
		Services:     []string{"example"},
		Dependencies: []string{"sqlite3-xorm"},
		Description:  "sqlite3-xorm example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
