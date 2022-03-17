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

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/erda-project/erda-infra/providers/mysql/v2/plugins/fields"
)

type User struct {
	ID        uint
	Name      string
	Age       uint
	DeletedAt fields.DeletedAt
}

func TestSoftDelete(t *testing.T) {
	DB, err := gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	DB = DB.Debug()
	if err != nil {
		t.Errorf("failed to connect database")
	}

	user := User{Name: "jinzhu", Age: 20}
	DB.Migrator().DropTable(&User{})
	DB.AutoMigrate(&User{})
	DB.Save(&user)

	var count int64
	var age uint

	if DB.Model(&User{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 1 {
		t.Errorf("Count soft deleted record, expects: %v, got: %v", 1, count)
	}

	if DB.Model(&User{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error != nil || age != user.Age {
		t.Errorf("Age soft deleted record, expects: %v, got: %v", 0, age)
	}

	if err := DB.Delete(&user).Error; err != nil {
		t.Fatalf("No error should happen when soft delete user, but got %v", err)
	}

	if user.DeletedAt.Time.IsZero() {
		t.Errorf("user's deleted at should not be zero, DeletedAt: %v", user.DeletedAt)
	}

	sql := DB.Session(&gorm.Session{DryRun: true}).Delete(&user).Statement.SQL.String()
	if !regexp.MustCompile(`UPDATE .users. SET .deleted_at.=.* WHERE .users.\..id. = .* AND \(.users.\..deleted_at. = \? OR .users.\..deleted_at. IS NULL\)`).MatchString(sql) {
		t.Fatalf("invalid sql generated, got %v", sql)
	}

	if DB.First(&User{}, "name = ?", user.Name).Error == nil {
		t.Errorf("Can't find a soft deleted record")
	}

	count = 0
	if DB.Model(&User{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 0 {
		t.Errorf("Count soft deleted record, expects: %v, got: %v", 0, count)
	}

	age = 0
	if err := DB.Model(&User{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error; err != nil || age != 0 {
		t.Fatalf("Age soft deleted record, expects: %v, got: %v, err %v", 0, age, err)
	}

	if err := DB.Unscoped().First(&User{}, "name = ?", user.Name).Error; err != nil {
		t.Errorf("Should find soft deleted record with Unscoped, but got err %s", err)
	}

	count = 0
	if DB.Unscoped().Model(&User{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 1 {
		t.Errorf("Count soft deleted record, expects: %v, count: %v", 1, count)
	}

	age = 0
	if DB.Unscoped().Model(&User{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error != nil || age != user.Age {
		t.Errorf("Age soft deleted record, expects: %v, got: %v", 0, age)
	}

	DB.Unscoped().Delete(&user)
	if err := DB.Unscoped().First(&User{}, "name = ?", user.Name).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Can't find permanently deleted record")
	}
}

type MilliUser struct {
	ID        uint
	Name      string
	Age       uint
	DeletedAt fields.DeletedAt `gorm:"softDelete:milli"`
}

func TestSoftDeleteMilliMode(t *testing.T) {
	DB, err := gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	DB = DB.Debug()
	if err != nil {
		t.Errorf("failed to connect database")
	}

	user := MilliUser{Name: "jinzhu", Age: 20}
	DB.Migrator().DropTable(&MilliUser{})
	DB.AutoMigrate(&MilliUser{})
	DB.Save(&user)

	var count int64
	var age uint

	if DB.Model(&MilliUser{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 1 {
		t.Errorf("Count soft deleted record, expects: %v, got: %v", 1, count)
	}

	if DB.Model(&MilliUser{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error != nil || age != user.Age {
		t.Errorf("Age soft deleted record, expects: %v, got: %v", 0, age)
	}

	if err := DB.Delete(&user).Error; err != nil {
		t.Fatalf("No error should happen when soft delete user, but got %v", err)
	}

	if user.DeletedAt.Time.IsZero() {
		t.Errorf("user's deleted at should not be zero, DeletedAt: %v", user.DeletedAt)
	}

	sql := DB.Session(&gorm.Session{DryRun: true}).Delete(&user).Statement.SQL.String()
	if !regexp.MustCompile(`UPDATE .milli_users. SET .deleted_at.=.* WHERE .milli_users.\..id. = .* AND \(.milli_users.\..deleted_at. = \? OR .milli_users.\..deleted_at. IS NULL\)`).MatchString(sql) {
		t.Fatalf("invalid sql generated, got %v", sql)
	}

	if DB.First(&MilliUser{}, "name = ?", user.Name).Error == nil {
		t.Errorf("Can't find a soft deleted record")
	}

	count = 0
	if DB.Model(&MilliUser{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 0 {
		t.Errorf("Count soft deleted record, expects: %v, got: %v", 0, count)
	}

	age = 0
	if err := DB.Model(&MilliUser{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error; err != nil || age != 0 {
		t.Fatalf("Age soft deleted record, expects: %v, got: %v, err %v", 0, age, err)
	}

	if err := DB.Unscoped().First(&MilliUser{}, "name = ?", user.Name).Error; err != nil {
		t.Errorf("Should find soft deleted record with Unscoped, but got err %s", err)
	}

	count = 0
	if DB.Unscoped().Model(&MilliUser{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 1 {
		t.Errorf("Count soft deleted record, expects: %v, count: %v", 1, count)
	}

	age = 0
	if DB.Unscoped().Model(&MilliUser{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error != nil || age != user.Age {
		t.Errorf("Age soft deleted record, expects: %v, got: %v", 0, age)
	}

	DB.Unscoped().Delete(&user)
	if err := DB.Unscoped().First(&MilliUser{}, "name = ?", user.Name).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Can't find permanently deleted record")
	}
}

type FlagUser struct {
	ID    uint
	Name  string
	Age   uint
	IsDel fields.DeletedAt `gorm:"softDelete:flag"`
}

func TestSoftDeleteFlagMode(t *testing.T) {
	DB, err := gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	DB = DB.Debug()
	if err != nil {
		t.Errorf("failed to connect database")
	}

	user := FlagUser{Name: "jinzhu", Age: 20}
	DB.Migrator().DropTable(&FlagUser{})
	DB.AutoMigrate(&FlagUser{})
	DB.Save(&user)

	var count int64
	var age uint

	if DB.Model(&FlagUser{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 1 {
		t.Errorf("Count soft deleted record, expects: %v, got: %v", 1, count)
	}

	if DB.Model(&FlagUser{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error != nil || age != user.Age {
		t.Errorf("Age soft deleted record, expects: %v, got: %v", 0, age)
	}

	if err := DB.Delete(&user).Error; err != nil {
		t.Fatalf("No error should happen when soft delete user, but got %v", err)
	}

	if user.IsDel.Time.IsZero() {
		t.Errorf("user's deleted at should not be zero, IsDel: %v", user.IsDel)
	}

	sql := DB.Session(&gorm.Session{DryRun: true}).Delete(&user).Statement.SQL.String()
	if !regexp.MustCompile(`UPDATE .flag_users. SET .is_del.=.* WHERE .flag_users.\..id. = .* AND \(.flag_users.\..is_del. = \? OR .flag_users.\..is_del. IS NULL\)`).MatchString(sql) {
		t.Fatalf("invalid sql generated, got %v", sql)
	}

	if DB.First(&FlagUser{}, "name = ?", user.Name).Error == nil {
		t.Errorf("Can't find a soft deleted record")
	}

	count = 0
	if DB.Model(&FlagUser{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 0 {
		t.Errorf("Count soft deleted record, expects: %v, got: %v", 0, count)
	}

	age = 0
	if err := DB.Model(&FlagUser{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error; err != nil || age != 0 {
		t.Fatalf("Age soft deleted record, expects: %v, got: %v, err %v", 0, age, err)
	}

	if err := DB.Unscoped().First(&FlagUser{}, "name = ?", user.Name).Error; err != nil {
		t.Errorf("Should find soft deleted record with Unscoped, but got err %s", err)
	}

	count = 0
	if DB.Unscoped().Model(&FlagUser{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 1 {
		t.Errorf("Count soft deleted record, expects: %v, count: %v", 1, count)
	}

	age = 0
	if DB.Unscoped().Model(&FlagUser{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error != nil || age != user.Age {
		t.Errorf("Age soft deleted record, expects: %v, got: %v", 0, age)
	}

	DB.Unscoped().Delete(&user)
	if err := DB.Unscoped().First(&FlagUser{}, "name = ?", user.Name).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Can't find permanently deleted record")
	}
}

type NullableDeletedAtUser struct {
	ID        int64
	Name      string
	Age       uint
	DeletedAt fields.DeletedAt `gorm:"default:null"`
}

func TestNullableDeletedAtUser(t *testing.T) {
	DB, err := gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	DB = DB.Debug()
	if err != nil {
		t.Errorf("failed to connect database")
	}

	user := NullableDeletedAtUser{Name: "shyamin", Age: 25}
	DB.Migrator().DropTable(&NullableDeletedAtUser{})
	DB.AutoMigrate(&NullableDeletedAtUser{})
	DB.Save(&user)

	var count int64
	var age uint

	if DB.Model(&NullableDeletedAtUser{}).Where("name = ?", user.Name).Count(&count).Error != nil || count != 1 {
		t.Errorf("Count soft deleted record, expects: %v, got: %v", 1, count)
	}

	if DB.Model(&NullableDeletedAtUser{}).Select("age").Where("name = ?", user.Name).Scan(&age).Error != nil || age != user.Age {
		t.Errorf("Age soft deleted record, expects: %v, got: %v", 0, age)
	}

	if err := DB.Delete(&user).Error; err != nil {
		t.Fatalf("No error should happen when soft delete user, but got %v", err)
	}

}

func TestSoftDeleteCreateClause_ModifyStatement(t *testing.T) {
	DB, err := gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	if err != nil {
		t.Errorf("failed to connect database")
	}
	DB = DB.Debug()

	user := User{Name: "dspo", Age: 20}
	DB.Migrator().DropTable(&User{})
	DB.AutoMigrate(&User{})
	DB.Save(&user)
	sql := DB.Session(&gorm.Session{DryRun: true}).Create(&user).Statement.SQL.String()
	t.Log(sql)
	sql = DB.Session(&gorm.Session{DryRun: true}).Delete(&user).Statement.SQL.String()
	t.Log(sql)
	DB.Delete(&user)
}
