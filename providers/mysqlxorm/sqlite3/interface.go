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
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"

	"github.com/erda-project/erda-infra/providers/mysqlxorm"
)

type Sqlite3 struct {
	db         *xorm.Engine
	closeState int32
}

func (s *Sqlite3) DB() *xorm.Engine {
	return s.db
}

func (s *Sqlite3) DataSourceName() string {
	return s.DB().DataSourceName()
}

func (s *Sqlite3) GetCloseState() bool {
	return atomic.LoadInt32(&s.closeState) == 1
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

	if o.RandomName {
		dbSourceName = randomName(dbSourceName)
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

	sqlite3Engine := &Sqlite3{
		db: engine,
	}

	return sqlite3Engine, nil
}

func (s *Sqlite3) Close() error {
	if atomic.CompareAndSwapInt32(&s.closeState, 0, 1) {
		err := s.DB().Close()
		if err != nil {
			return err
		}

		err = os.Remove(s.DataSourceName())
		return err
	}
	return nil
}

func randomName(path string) string {
	dir, file := filepath.Split(path)
	name := strings.TrimSuffix(file, filepath.Ext(file))
	random := fmt.Sprintf("%s-%s%s", name, strings.ReplaceAll(uuid.New().String(), "-", ""), filepath.Ext(file))
	return filepath.Join(dir, random)
}
