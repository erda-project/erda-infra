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

	_ "github.com/mattn/go-sqlite3"
	"github.com/xormplus/xorm"

	"github.com/erda-project/erda-infra/providers/mysqlxorm"
)

type Sqlite struct {
	db *xorm.Engine
}

func (s *Sqlite) DB() *xorm.Engine {
	return s.db
}

func (s *Sqlite) NewSession(ops ...mysqlxorm.SessionOption) *mysqlxorm.Session {
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
func NewSqlite3(dbSourceName string) (*Sqlite, error) {
	if dbSourceName == "" {
		return nil, errors.New("empty dbSourceName")
	}

	sqlite3, err := xorm.NewSqlite3(dbSourceName)

	if err != nil {
		return nil, err
	}

	sqlite3Engine := &Sqlite{db: sqlite3}

	return sqlite3Engine, nil
}
