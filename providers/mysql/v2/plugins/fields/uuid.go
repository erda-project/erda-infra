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

package fields

import (
	"database/sql"
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// UUID auto generate an uuid on creating record.
type UUID sql.NullString

// Scan .
func (u *UUID) Scan(value interface{}) error {
	return (*sql.NullString)(u).Scan(value)
}

// Value .
func (u UUID) Value() (driver.Value, error) {
	return sql.NullString(u).Value()
}

// MustValue returns u.String
func (u UUID) MustValue() string {
	return u.String
}

// CreateClauses .
func (UUID) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{UUIDCreateClause{Field: f}}
}

// UUIDCreateClause .
type UUIDCreateClause struct {
	Field *schema.Field
}

// Name .
func (cc UUIDCreateClause) Name() string {
	return ""
}

// Build .
func (cc UUIDCreateClause) Build(clause.Builder) {
}

// MergeClause .
func (cc UUIDCreateClause) MergeClause(*clause.Clause) {
}

// ModifyStatement .
func (cc UUIDCreateClause) ModifyStatement(stmt *gorm.Statement) {
	stmt.SetColumn(cc.Field.Name, UUID{
		String: uuid.New().String(),
		Valid:  true,
	})
}
