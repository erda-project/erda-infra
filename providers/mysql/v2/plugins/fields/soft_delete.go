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
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var zero = time.Unix(0, 0)

// DeletedAt on delete: the field value will be set to the current datetime.
// on query: the record will be not returned if it is soft-deleted.
// if you want to find the record which is soft-deleted, or you want
// to delete the record forever, you can use db.Unscoped().Find(&users),
// and db.Unscoped().Delete(&order).
type DeletedAt sql.NullTime

// Scan implements the Scanner interface.
func (n *DeletedAt) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n DeletedAt) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

// MustValue returns zero timestamp if the value is null;
// returns time value if the value is not null.
// if the value is null returns
func (n DeletedAt) MustValue() time.Time {
	if n.Valid {
		return n.Time
	}
	return zero
}

// MarshalJSON .
func (n DeletedAt) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *DeletedAt) UnmarshalJSON(b []byte) error {
	bs := string(b)
	if strings.EqualFold(bs, "null") {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Time)
	n.Valid = err == nil
	return err
}

// QueryClauses .
func (DeletedAt) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteQueryClause{Field: f}}
}

// SoftDeleteQueryClause .
type SoftDeleteQueryClause struct {
	Field *schema.Field
}

// Name .
func (sd SoftDeleteQueryClause) Name() string {
	return ""
}

// Build .
func (sd SoftDeleteQueryClause) Build(clause.Builder) {
}

// MergeClause .
func (sd SoftDeleteQueryClause) MergeClause(*clause.Clause) {
}

// ModifyStatement .
func (sd SoftDeleteQueryClause) ModifyStatement(stmt *gorm.Statement) {
	if _, ok := stmt.Clauses["soft_delete_enabled"]; !ok && !stmt.Statement.Unscoped {
		if c, ok := stmt.Clauses["WHERE"]; ok {
			if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) >= 1 {
				for _, expr := range where.Exprs {
					if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
						where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
						c.Expression = where
						stmt.Clauses["WHERE"] = c
						break
					}
				}
			}
		}

		stmt.AddClause(clause.Where{Exprs: []clause.Expression{
			clause.Or(
				clause.Lte{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: zero},
				clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: nil},
			),
		}})
		stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
	}
}

// UpdateClauses .
func (DeletedAt) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteUpdateClause{Field: f}}
}

// SoftDeleteUpdateClause .
type SoftDeleteUpdateClause struct {
	Field *schema.Field
}

// Name .
func (sd SoftDeleteUpdateClause) Name() string {
	return ""
}

// Build .
func (sd SoftDeleteUpdateClause) Build(clause.Builder) {
}

// MergeClause .
func (sd SoftDeleteUpdateClause) MergeClause(*clause.Clause) {
}

// ModifyStatement .
func (sd SoftDeleteUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		SoftDeleteQueryClause(sd).ModifyStatement(stmt)
	}
}

// DeleteClauses .
func (DeletedAt) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteDeleteClause{Field: f}}
}

// SoftDeleteDeleteClause .
type SoftDeleteDeleteClause struct {
	Field *schema.Field
}

// Name .
func (sd SoftDeleteDeleteClause) Name() string {
	return ""
}

// Build .
func (sd SoftDeleteDeleteClause) Build(clause.Builder) {
}

// MergeClause .
func (sd SoftDeleteDeleteClause) MergeClause(*clause.Clause) {
}

// ModifyStatement .
func (sd SoftDeleteDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		curTime := stmt.DB.NowFunc()
		stmt.AddClause(clause.Set{{Column: clause.Column{Name: sd.Field.DBName}, Value: curTime}})
		stmt.SetColumn(sd.Field.DBName, curTime, true)

		if stmt.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.Context, stmt.ReflectValue, stmt.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}

			if stmt.ReflectValue.CanAddr() && stmt.Dest != stmt.Model && stmt.Model != nil {
				_, queryValues = schema.GetIdentityFieldValuesMap(stmt.Context, reflect.ValueOf(stmt.Model), stmt.Schema.PrimaryFields)
				column, values = schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}
			}
		}

		SoftDeleteQueryClause(sd).ModifyStatement(stmt)
		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.Build(stmt.DB.Callback().Update().Clauses...)
	}
}

// CreateClauses .
func (DeletedAt) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{SoftDeleteCreateClause{Field: f}}
}

// SoftDeleteCreateClause .
type SoftDeleteCreateClause struct {
	Field *schema.Field
}

// Name .
func (sd SoftDeleteCreateClause) Name() string {
	return ""
}

// Build .
func (sd SoftDeleteCreateClause) Build(clause.Builder) {
}

// MergeClause .
func (sd SoftDeleteCreateClause) MergeClause(*clause.Clause) {
}

// ModifyStatement .
func (sd SoftDeleteCreateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		stmt.SetColumn(sd.Field.Name, DeletedAt{Time: zero, Valid: true})
	}
}
