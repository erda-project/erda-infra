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
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/erda-project/erda-infra/providers/mysqlxorm"
)

type Sqlite3 struct {
	db *xorm.Engine
}

func (s *Sqlite3) DB() *xorm.Engine {
	return s.db
}

func (s *Sqlite3) NewSession(ops ...mysqlxorm.SessionOption) *mysqlxorm.Session {
	tx := &mysqlxorm.Session{}
	for _, opt := range ops {
		opt(tx)
	}

	if tx.Session == nil {
		tx.Session = s.db.NewSession()
	}

	return tx
}

// NewSqlite3 Use for unit-test
func NewSqlite3(dbSourceName string, opts ...OptionFunc) (*Sqlite3, error) {
	if dbSourceName == "" {
		return nil, errors.New("empty dbSourceName")
	}

	o := &Options{}

	for _, opt := range opts {
		opt(o)
	}

	engine, err := xorm.NewEngine("sqlite3", dbSourceName)
	if err != nil {
		return nil, err
	}

	// set journal_mode in sqlite3
	// the default journal_mode in sqlite is `delete`
	if o.JournalMode != "" {
		_, err = engine.Exec(fmt.Sprintf("PRAGMA journal_mode = %s", o.JournalMode))
		if err != nil {
			return nil, err
		}
	}

	engine.SetMapper(names.GonicMapper{})

	sqlite3Engine := &Sqlite3{db: engine}

	return sqlite3Engine, nil
}
