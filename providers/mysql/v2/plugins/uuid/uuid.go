// Copyright (c) 2022 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package uuid

import (
	"database/sql"
	"database/sql/driver"

	googleUUID "github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type UUID sql.NullString

func (u *UUID) Scan(value interface{}) error {
	return (*sql.NullString)(u).Scan(value)
}

func (u UUID) Value() (driver.Value, error) {
	return sql.NullString(u).Value()
}

//func (u UUID) MarshalJSON() ([]byte, error) {
//	if u.Valid {
//		return json.Marshal(u.String)
//	}
//	return json.Marshal(nil)
//}
//
//func (u *UUID) UnmarshalJSON(b []byte) error {
//	if string(b) == "null" {
//		u.Valid = false
//		u.String = ""
//		return nil
//	}
//
//}

func (UUID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{UUIDCreateClause{Field: f}}
}

type UUIDCreateClause struct {
	Field *schema.Field
}

func (sd UUIDCreateClause) Name() string {
	return ""
}

func (sd UUIDCreateClause) Build(clause.Builder) {
}

func (sd UUIDCreateClause) MergeClause(*clause.Clause) {
}

func (sd UUIDCreateClause) ModifyStatement(stmt *gorm.Statement) {
	stmt.SetColumn(sd.Field.Name, UUID{
		String: googleUUID.New().String(),
		Valid:  true,
	})
}
