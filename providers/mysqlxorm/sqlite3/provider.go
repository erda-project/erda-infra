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

package sqlite3

import (
	"fmt"
	"reflect"

	"github.com/xormplus/xorm"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/mysqlxorm"
)

type config struct {
	DbSourceName string `file:"dbSourceName" env:"DB_SOURCE_NAME" default:"test.sqlite3"`
}

type provider struct {
	Cfg *config
	Log logs.Logger
	*Sqlite
}

var _ servicehub.ProviderInitializer = (*provider)(nil)

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	server, err := xorm.NewSqlite3(p.Cfg.DbSourceName)
	if err != nil {
		return fmt.Errorf("failed to connect to sqlite3 server,err : %s", err)
	}

	p.Sqlite = &Sqlite{db: server}

	return nil
}

func init() {
	servicehub.Register("sqlite3-xorm", &servicehub.Spec{
		Services: []string{"sqlite3-xorm"},
		Types: []reflect.Type{
			reflect.TypeOf((*mysqlxorm.Interface)(nil)).Elem(),
		},
		Description: "sqlite3-xorm",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
