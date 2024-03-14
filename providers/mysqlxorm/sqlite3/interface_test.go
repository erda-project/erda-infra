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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"xorm.io/xorm/names"

	"github.com/erda-project/erda-infra/providers/mysqlxorm"
)

const dbSourceName = "test1.db"

type Server struct {
	mysql mysqlxorm.Interface
}

type User struct {
	ID   uint64 `json:"id" xorm:"pk autoincr"`
	Name string `json:"name"`
}

func (u *User) TableName() string {
	return "user"
}

func (s *Server) GetUserByID(id uint64, ops ...mysqlxorm.SessionOption) (*User, error) {
	session := s.mysql.NewSession(ops...)
	defer session.Close()

	var user User
	_, err := session.ID(id).Get(&user)

	return &user, err
}

func (s *Server) CreateUser(user *User, ops ...mysqlxorm.SessionOption) (err error) {
	session := s.mysql.NewSession(ops...)
	defer session.Close()

	_, err = session.Insert(user)
	return err
}

func (s *Server) TestTx(err error, ops ...mysqlxorm.SessionOption) error {
	session := s.mysql.NewSession(ops...)
	defer session.Close()
	return err
}

func TestNewSqlite3(t *testing.T) {
	dbname := filepath.Join(os.TempDir(), dbSourceName)
	engine, err := NewSqlite3(dbname)
	if err != nil {
		t.Fatalf("new sqlite3 err : %s", err)
	}

	defer engine.Close()

	server := Server{
		mysql: engine,
	}

	server.mysql.DB().SetMapper(names.GonicMapper{})
	server.mysql.DB().Sync2(&User{})

	testCase := []struct {
		name       string
		insertUser []User
		txErr      error
		want       []User
	}{
		{
			name:  "test tx",
			txErr: errors.New("tx error"),
			insertUser: []User{
				{ID: 4, Name: "Alice"},
				{ID: 5, Name: "Bob"},
				{ID: 6, Name: "Cat"},
			},
			want: []User{},
		},
		{
			name: "sqlite3 use for xorm",
			insertUser: []User{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
				{ID: 3, Name: "Cat"},
			},
			txErr: nil,
			want: []User{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
				{ID: 3, Name: "Cat"},
			},
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			tx := server.mysql.NewSession()
			defer tx.Close()
			if err = tx.Begin(); err != nil {
				t.Fatalf("tx begin err : %s", err)
			}

			ops := mysqlxorm.WithSession(tx)
			// insert sql
			for _, user := range test.insertUser {
				err = server.CreateUser(&user, ops)
				if err != nil {
					tx.Rollback()
					t.Fatalf("create user err : %s", err)
				}
			}

			err = server.TestTx(test.txErr, ops)
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}

			if len(test.want) <= 0 {
				for _, user := range test.insertUser {
					u, err := server.GetUserByID(user.ID)
					if err != nil {
						t.Fatalf("get user err : %s", err)
					}
					assert.Equal(t, &User{}, u)
				}
				return
			}

			for _, user := range test.want {
				u, err := server.GetUserByID(user.ID)
				if err != nil {
					t.Fatalf("get user err : %s", err)
				}
				assert.Equal(t, user.Name, u.Name)
			}
		})
	}
}

func TestJournalMode(t *testing.T) {
	dbname := filepath.Join(os.TempDir(), dbSourceName)

	want := []JournalMode{
		MEMORY,
		DELETE,
		PERSIST,
		OFF,
		WAL,
		TRUNCATE,
	}
	defer func() {
		os.Remove(dbname)
	}()

	for _, w := range want {
		engine, err := NewSqlite3(dbname, WithJournalMode(w))
		if err != nil {
			t.Fatalf("new sqlite3 err : %s", err)
		}
		defer engine.Close()

		// get journal in sqlite3
		results, _ := engine.DB().Query("PRAGMA journal_mode;")
		assert.Equal(t, string(w), string(results[0]["journal_mode"]))
		engine.Close()
	}

}

func TestRandomName(t *testing.T) {
	path := "/test/dir/sample.txt"
	name1 := randomName(path)
	name2 := randomName(path)

	assert.True(t, strings.HasPrefix(name1, "/test/dir/sample-"), "Random name does not start with original name")

	assert.Equal(t, filepath.Ext(name1), ".txt", "Random name does not have original extension")

	assert.NotEqual(t, name1, name2, "Random name generator produced the same result twice")
}

func TestWithRandomName(t *testing.T) {
	dbname := filepath.Join(os.TempDir(), dbSourceName)
	engine, err := NewSqlite3(dbname, WithRandomName(false))
	if err != nil {
		panic(err)
	}
	defer engine.Close()
	assert.Nil(t, err)
	assert.Equal(t, dbname, engine.DataSourceName())
	engine.Close()

	engineRandom, err := NewSqlite3(dbname, WithRandomName(true))
	assert.Nil(t, err)
	assert.NotEqual(t, dbname, engineRandom, "Random name is not take effect")
	assert.Equal(t, filepath.Ext(engineRandom.DataSourceName()), filepath.Ext(dbname), "Random names does not have original extension")
	defer engineRandom.Close()
}
