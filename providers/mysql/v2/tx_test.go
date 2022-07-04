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

package v2_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	v2 "github.com/erda-project/erda-infra/providers/mysql/v2"
)

var dsn = filepath.Join(os.TempDir(), "gorm.db")
var tx *v2.TX

type User struct {
	Age  int64 `gorm:"type:BIGINT"`
	Name string
}

func TestTX(t *testing.T) {
	openDB(t)
	defer closeDB()

	if err := tx.DB().AutoMigrate(new(User)); err != nil {
		t.Fatalf("failed to migrate user: %v", err)
	}

	var (
		name = v2.Col("name")
		age  = v2.Col("age")
	)

	var prepare = func(t *testing.T) {
		var users []User
		for i := 0; i < 10; i++ {
			users = append(users, User{Age: int64(i), Name: "dspo-" + strconv.Itoa(i)})
		}
		// INSERT INTO `users` (`age`,`name`) VALUES (0,"dspo-0"),(1,"dspo-1"),(2,"dspo-2"),(3,"dspo-3"),(4,"dspo-4"),(5,"dspo-5"),(6,"dspo-6"),(7,"dspo-7"),(8,"dspo-8"),(9,"dspo-9")
		if err := tx.CreateInBatches(users, len(users)); err != nil {
			t.Fatal(err)
		}
		total, err := tx.List(new([]User))
		if err != nil {
			t.Fatal(err)
		}
		if total != 10 {
			t.Fatalf("expects total: %v, got: %v", 11, total)
		}
	}
	var clear = func(t *testing.T) {
		// DELETE FROM `users` WHERE 1=1
		_, err := tx.Delete(new(User), v2.Where("1=1"))
		if err != nil {
			t.Fatal(err)
		}
		total, err := tx.List(new([]User))
		if err != nil {
			t.Fatal(err)
		}
		if total != 0 {
			t.Fatalf("expects total: %v, got: %v", 0, total)
		}
	}

	t.Run("TX.Create", func(t *testing.T) {
		defer t.Run("clear", clear)

		// INSERT INTO `users` (`age`,`name`) VALUES (10,"dspo-10")
		if err := tx.Create(&User{Age: 10, Name: "dspo-10"}); err != nil {
			t.Fatalf("failed to Create: %v", err)
		}
		var user User
		// SELECT * FROM `users` WHERE age = 10 AND name = "dspo-10" ORDER BY `users`.`age` LIMIT 1
		ok, err := tx.Get(&user, age.Is(10), name.Is("dspo-10"))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("not ok")
		}
	})

	t.Run("TX.CreateInBatches", func(t *testing.T) {
		prepare(t)
		clear(t)
	})

	t.Run("TX.Delete", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		for i := 20; i < 30; i++ {
			users = append(users, User{Age: int64(i), Name: "six-" + strconv.Itoa(i)})
		}
		if err := tx.CreateInBatches(users, len(users)); err != nil {
			t.Fatal(err)
		}
		total, err := tx.Delete(new(User), age.GreaterThan(30))
		if err != nil {
			t.Fatal(err)
		}
		if total != 0 {
			t.Fatalf("expects total: %v, got: %v", 0, total)
		}

		total, err = tx.Delete(new(User), name.Like("six-%"))
		if err != nil {
			t.Fatal(err)
		}
		if total != 10 {
			t.Fatalf("expects total: %v, got: %v", 10, total)
		}
	})

	t.Run("TX.Updates by map", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var user User
		var values = map[string]interface{}{
			"name": "cmc-10",
			"age":  int64(100),
		}
		if err := tx.Updates(&user, values, name.Is("dspo-10")); err != nil {
			t.Fatal(err)
		}
		if user.Name != values["name"] {
			t.Fatalf("expects user.Name: %v, got: %v", values["name"], user.Name)
		}
		if user.Age != values["age"].(int64) {
			t.Fatalf("expects user.Age: %v, got: %v", values["age"], user.Age)
		}
	})

	t.Run("TX.Updates by struct", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var user User
		var values = User{
			Age:  0,
			Name: "0",
		}
		if err := tx.Updates(&user, values, name.Is("cmc-10")); err != nil {
			t.Fatal(err)
		}
		if user.Age != values.Age {
			t.Fatalf("expects user.Age: %v, got: %v", values, user.Age)
		}
		if user.Name != values.Name {
			t.Fatalf("expects user.Name: %v, got: %v", values.Name, user.Name)
		}
	})

	t.Run("TX.SetColumns", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var user User
		var value = User{Name: "cmc-10", Age: 10}
		if err := tx.SetColumns(&user, name.Is("0"), name.Set(value.Name), age.Set(value.Age)); err != nil {
			t.Fatal(err)
		}
		if user.Name != value.Name {
			t.Fatalf("expects user.Name: %v, got: %v", value.Name, user.Name)
		}
		if user.Age != value.Age {
			t.Fatalf("expects user.Age: %v, got: %v", value.Age, user.Age)
		}
	})

	t.Run("Where", func(t *testing.T) {
		t.Log("ignore")
	})

	t.Run("Wheres", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var user User
		ok, err := tx.Get(&user, v2.Wheres(map[string]interface{}{"name": "dspo-1"}))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("expects ok: %v, got: %v", true, ok)
		}
		if user.Name != "dspo-1" || user.Age != 1 {
			t.Fatalf("expects user.Name: %v, user.Age: %v, got name: %v, age: %v",
				"dspo-1", 1, user.Name, user.Age)
		}
	})

	t.Run("TX.List", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users)
		if err != nil {
			t.Fatal(err)
		}
		if total != 10 {
			t.Fatalf("expects total: %v, got: %v", 11, total)
		}
		if len(users) != 10 {
			t.Fatalf("expects length of users: %v, got: %v", 11, len(users))
		}
	})

	t.Run("TX.List with conditions", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, name.Like("dspo%"), age.GreaterThan(5))
		if err != nil {
			t.Fatal(err)
		}
		if total != 4 {
			t.Fatalf("expects total: %v, got: %v", 4, total)
		}
		if len(users) != 4 {
			t.Fatalf("expects length of users: %v, got: %v", 4, len(users))
		}
	})

	t.Run("TX.List order by DESC", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, age.DESC())
		if err != nil {
			t.Fatal(err)
		}
		if total != 10 {
			t.Fatalf("expects total: %v, got: %v", 11, total)
		}
		if users[0].Age != 9 {
			t.Fatalf("expects users[0].Age: %v, got: %v", 10, users[0].Age)
		}
	})

	t.Run("TX.List order by ASC", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, age.ASC())
		if err != nil {
			t.Fatal(err)
		}
		if total != 10 {
			t.Fatalf("expects total: %v, got: %v", 11, total)
		}
		if users[0].Age != 0 {
			t.Fatalf("expects users[0].Age: %v, got: %v", 0, users[0].Age)
		}
	})

	t.Run("TX.List by paging", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, v2.Paging(5, 1))
		if err != nil {
			t.Fatal(err)
		}
		if total != 10 {
			t.Fatalf("expects total: %v, got: %v", 11, total)
		}
		if len(users) != 5 {
			t.Fatalf("expects length of users: %v, got: %v", 5, len(users))
		}

		if total != 10 {
			t.Fatalf("expects total: %v, got: %v", 11, total)
		}
		if len(users) != 5 {
			t.Fatalf("expects length of users: %v, got: %v", 5, len(users))
		}

		// SELECT * FROM `users`
		total, err = tx.List(&users, v2.Paging(-1, 0))
		if err != nil {
			t.Fatal(err)
		}
		if total != 10 {
			t.Fatalf("expects total: %v, got: %v", 10, total)
		}
	})

	t.Run("TX.Get", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var user User
		ok, err := tx.Get(&user)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("expects ok: %v, got: %v", true, ok)
		}
	})

	t.Run("Tx.Get with conditions", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var user User
		ok, err := tx.Get(&user, name.Is("dspo-2"))
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("expects ok: %v, got: %v", true, ok)
		}

		ok, err = tx.Get(&user, name.Is("xxx"))
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Fatalf("expects ok: %v, got: %v", false, ok)
		}
	})

	t.Run("WhereColumn.Is", func(t *testing.T) {
		t.Log("ignore")
	})

	t.Run("WhereColumn.In", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, name.In([]interface{}{"dspo-1", "dspo-2", "six-1", "six-2"}))
		if err != nil {
			t.Fatal(err)
		}
		if total != 2 {
			t.Fatalf("expects total: %v, got: %v", 2, total)
		}
		if len(users) != 2 {
			t.Fatalf("expects length of users: %v, got: %v", 2, len(users))
		}
	})

	t.Run("WhereColumn.InMap", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, name.InMap(map[interface{}]struct{}{
			"dspo-1": {},
			"dspo-2": {},
			"six-1":  {},
			"six-2":  {},
		}))
		if err != nil {
			t.Fatal(err)
		}
		if total != 2 {
			t.Fatalf("expects total: %v, got: %v", 2, total)
		}
		if len(users) != 2 {
			t.Fatalf("expects length of users: %v, got: %v", 2, len(users))
		}
	})

	t.Run("WhereColumn.Like", func(t *testing.T) {
		t.Log("ignore")
	})

	t.Run("WhereColumn.GreaterThan", func(t *testing.T) {
		t.Log("ignore")
	})

	t.Run("WhereColumn.EqGreaterThan", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, age.EqGreaterThan(5))
		if err != nil {
			t.Fatal(err)
		}
		if total != 5 {
			t.Fatalf("expects total: %v, got: %v", 5, total)
		}
	})

	t.Run("WhereColumn.LessThan", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, age.LessThan(5))
		if err != nil {
			t.Fatal(err)
		}
		if total != 5 {
			t.Fatalf("expects total: %v, got: %v", 5, total)
		}
	})

	t.Run("WhereColumn.EqLessThan", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, age.EqLessThan(5))
		if err != nil {
			t.Fatal(err)
		}
		if total != 6 {
			t.Fatalf("expects total: %v, got: %v", 6, total)
		}
	})

	t.Run("OrderColumn.DESC", func(t *testing.T) {
		t.Log("ignore")
	})

	t.Run("OrderColumn.ASC", func(t *testing.T) {
		t.Log("ignore")
	})

	t.Run("SetColumn.Set", func(t *testing.T) {
		t.Log("ignore")
	})

	t.Run("WhereValue.In", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, v2.Value("dspo-1").In("name", "age"))
		if err != nil {
			t.Fatal(err)
		}
		if total != 1 {
			t.Fatalf("expects total: %v, got: %v", 1, total)
		}
	})

	t.Run("Paging", func(t *testing.T) {
		t.Log("ignore")
	})

	t.Run("OrderBy", func(t *testing.T) {
		prepare(t)
		defer clear(t)

		var users []User
		total, err := tx.List(&users, v2.OrderBy("age", v2.DESC))
		if err != nil {
			t.Fatal(err)
		}
		if total != 10 {
			t.Fatalf("expect total: %v, got: %v", 10, total)
		}
		if users[0].Name != "dspo-9" || users[0].Age != 9 {
			t.Fatalf("expects user[0].Name: %v, user[0].Age: %v, got name: %v, age: %v",
				"dspo-0", 0, users[0].Name, users[0].Age)
		}

		_, err = tx.List(&users, v2.OrderBy("age", v2.ASC))
		if err != nil {
			t.Fatal(err)
		}
		if users[0].Name != "dspo-0" || users[0].Age != 0 {
			t.Fatalf("expects user[0].Name: %v, user[0].Age: %v, got name: %v, age: %v",
				"dspo-0", 0, users[0].Name, users[0].Age)
		}
	})
}

func openDB(t *testing.T) {
	if tx != nil {
		return
	}
	closeDB()
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open %s: %v", dsn, err)
	}
	tx = v2.NewTx(db.Debug())
}

func closeDB() {
	os.Remove(dsn)
}
