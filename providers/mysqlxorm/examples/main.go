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

	"github.com/sirupsen/logrus"
	"xorm.io/xorm"

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
	/* test db */
	r, err := p.MySQL.DB().QueryString("show tables")
	if err != nil {
		panic(err)
	}
	for i, m := range r {
		for k, v := range m {
			fmt.Println(i, k, v)
		}
	}

	/* test tx */
	// create table for test
	if err := p.DB.CreateTables(&Table{}); err != nil {
		return err
	}
	defer func() {
		if err := p.DB.DropTables(&Table{}); err != nil {
			logrus.Fatalf("failed to cleanup table, err: %v", err)
		}
	}()

	// tx
	txSession := p.MySQL.NewSession()
	defer txSession.Close()
	// begin tx
	err = txSession.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			logrus.Error(err)

			// check table size before rollback, should be 0 (2 in tx)
			count := p.checkTableRows(0)
			logrus.Printf("count before rollback: %d", count)

			// rollback
			err = txSession.Rollback()
			if err != nil {
				logrus.Fatalf("failed to rollback, err: %v", err)
			}
		} else { // make insertFailed as success
			p.checkTableRows(2)
		}
	}()

	// insert success
	err = p.insertSuccess(mysqlxorm.WithSession(txSession))
	if err != nil {
		return err
	}

	// insert failed
	err = p.insertFailed(mysqlxorm.WithSession(txSession))
	if err != nil {
		return err
	}

	// tx commit
	err = txSession.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Table .
type Table struct {
	ID   uint64 `json:"id" xorm:"pk autoincr"`
	Name string
}

// TableName .
func (t *Table) TableName() string { return "table" }

func (p *provider) insertSuccess(opts ...mysqlxorm.SessionOption) error {
	s := p.MySQL.NewSession(opts...)
	defer s.Close()
	_, err := s.InsertOne(&Table{Name: "n1"})
	return err
}

func (p *provider) insertFailed(opts ...mysqlxorm.SessionOption) error {
	s := p.MySQL.NewSession(opts...)
	defer s.Close()
	_, err := s.InsertOne(&Table{Name: "n2"})
	if err != nil {
		return err
	}
	// force err
	return fmt.Errorf("fake error by func: insertFailed")
}

func (p *provider) checkTableRows(expectRows int64) int64 {
	count, _ := p.DB.Count(&Table{})
	if count != expectRows {
		logrus.Errorf("expectRows: %d, actualRows: %d", expectRows, count)
	} else {
		logrus.Infof("expectRows: %d, actualRows: %d", expectRows, count)
	}
	return count
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
