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
	"os"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/erda-project/erda-infra/providers/mysql/v2/plugins/fields"
)

type UUIDUser struct {
	ID   fields.UUID
	Name string
}

func TestCreateClause_ModifyStatement(t *testing.T) {
	DB, err := gorm.Open(sqlite.Open(filepath.Join(os.TempDir(), "gorm.db")), &gorm.Config{})
	DB = DB.Debug()
	if err != nil {
		t.Errorf("failed to connect database")
	}

	user := UUIDUser{Name: "dspo"}
	DB.Migrator().DropTable(&UUIDUser{})
	DB.AutoMigrate(&UUIDUser{})
	DB.Save(&user)
	sql := DB.Session(&gorm.Session{DryRun: true}).Save(&user).Statement.SQL.String()
	t.Log(sql)
}
