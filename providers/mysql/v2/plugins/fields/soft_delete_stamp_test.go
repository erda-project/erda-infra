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

package fields_test

//import (
//	"errors"
//	"os"
//	"path/filepath"
//	"regexp"
//	"testing"
//
//	"gorm.io/driver/sqlite"
//	"gorm.io/gorm"
//
//	"github.com/erda-project/erda-infra/providers/mysql/v2/plugins/fields"
//)
//
//type User2 struct {
//	ID        uint
//	Name      string
//	Age       uint
//	DeletedAt fields.DeletedAtStamp
//}
//
//func TestDeletedAtStamp(t *testing.T) {
//	dsn := filepath.Join(os.TempDir(), "gorm.db")
//	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
//	if err != nil {
//		t.Fatalf("failed to connect to database: %v", err)
//	}
//	defer os.Remove(dsn)
//	db = db.Debug()
//
//	user := User2{Name: "dspo", Age: 20}
//	t.Log("drop table")
//	if err = db.Migrator().DropTable(new(User2)); err != nil {
//		t.Fatalf("failed to drop table: %v", err)
//	}
//	t.Log("auto migrate")
//	if err = db.AutoMigrate(new(User2)); err != nil {
//		t.Fatalf("failed to auto migrate: %v", err)
//	}
//	t.Log("save user")
//	db.Save(&user)
//
//	var count int64
//	var age uint
//
//	t.Log("count")
//	if db.Model(&User2{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 1 {
//		t.Errorf("count soft deleted record, expects: %v, got: %v", 1, count)
//	}
//
//	t.Log("select")
//	if db.Model(&User2{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error != nil || age != user.Age {
//		t.Errorf("age soft deleted record, expects: %v, got: %v", user.Age, age)
//	}
//
//	t.Log("delete")
//	if err := db.Delete(&user).Error; err != nil {
//		t.Fatalf("no error should be happen when soft delte user, but got: %v", err)
//	}
//
//	if user.DeletedAt.Int64 == 0 {
//		t.Errorf("user's deleted at shoud be zero, bug got: %v", user.DeletedAt)
//	}
//
//	t.Log("dry run delete")
//	sql := db.Session(&gorm.Session{DryRun: true}).Delete(&user).Statement.SQL.String()
//	if !regexp.MustCompile(`UPDATE .user2. SET .deleted_at.=.* WHERE .user2.\..id. = .* AND \(.user2.\..deleted_at. = \? OR .user2.\..deleted_at. IS NULL\)`).MatchString(sql) {
//		t.Fatalf("invalid sql generated, got %v", sql)
//	}
//
//	t.Log("first")
//	if db.First(&User2{}, "name = ?", user.Name).Error == nil {
//		t.Errorf("can not find a soft deleted record")
//	}
//
//	count = 0
//	t.Log("count")
//	if db.Model(&User2{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 0 {
//		t.Errorf("count soft deleted record, expects: %v, got: %v", 0, count)
//	}
//
//	age = 0
//	t.Log("select age")
//	if err := db.Model(&User2{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error; err != nil || age != 0 {
//		t.Fatalf("age soft deleted record, expects: %v, got: %v, err: %v", 0, age, err)
//	}
//
//	t.Log("unscoped first")
//	if err := db.Unscoped().First(&User2{}, "name = ?", user.Name).Error; err != nil {
//		t.Errorf("should find soft deleted record with Unscoped, but got err: %v", err)
//	}
//
//	t.Log("unscoped delete")
//	db.Unscoped().Delete(&user)
//	t.Log("unscoped first")
//	if err := db.Unscoped().First(&User2{}, "name = ?", user.Name).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
//		t.Errorf("can not permanently deleted record")
//	}
//}
