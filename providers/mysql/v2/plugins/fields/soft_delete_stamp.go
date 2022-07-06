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

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// DeletedAtStamp is the field type for soft deleting with BIGINT type.
// 0 or null presents not deleted, an UnixMilli timestamp presents deleted.
type DeletedAtStamp sql.NullInt64

// Scan .
func (n *DeletedAtStamp) Scan(value interface{}) error {
	return (*sql.NullInt64)(n).Scan(value)
}

// Value .
func (n DeletedAtStamp) Value() (driver.Value, error) {
	return (sql.NullInt64)(n).Value()
}

// MarshalJSON .
func (n DeletedAtStamp) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON .
func (n *DeletedAtStamp) UnmarshalJSON(b []byte) error {
	bs := string(b)
	if strings.EqualFold(bs, "null") {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Int64)
	n.Valid = err == nil
	return err

}

// GormDataType .
func (DeletedAtStamp) GormDataType() string {
	return "BIGINT(20)"
}

// QueryClauses .
func (DeletedAtStamp) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{DeletedatstampQueryclause{Field: f}}
}

// UpdateClauses .
func (DeletedAtStamp) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{DeletedAtStampUpdateClause{Field: f}}
}

// DeleteClauses .
func (DeletedAtStamp) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{DeletedAtStampDeleteClause{Field: f}}
}

// CreateClauses .
func (DeletedAtStamp) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{DeletedAtStampCreateClause{Field: f}}
}

type baseClause struct {
	softDeletedMode string
}

// Name .
func (baseClause) Name() string {
	return ""
}

// Build .
func (baseClause) Build(_ clause.Builder) {}

// MergeClause .
func (baseClause) MergeClause(_ *clause.Clause) {}

// DeletedatstampQueryclause .
type DeletedatstampQueryclause struct {
	baseClause

	Field *schema.Field
}

// ModifyStatement .
func (qc DeletedatstampQueryclause) ModifyStatement(stmt *gorm.Statement) {
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
				clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: qc.Field.DBName}, Value: 0},
				clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: qc.Field.DBName}, Value: nil},
			),
		}})
		stmt.Clauses["soft_delete_enabled"] = clause.Clause{}
	}
}

// DeletedAtStampUpdateClause .
type DeletedAtStampUpdateClause struct {
	baseClause

	Field *schema.Field
}

// ModifyStatement .
func (qc DeletedAtStampUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		DeletedatstampQueryclause(qc).ModifyStatement(stmt)
	}
}

// DeletedAtStampDeleteClause .
type DeletedAtStampDeleteClause struct {
	baseClause

	Field *schema.Field
}

// ModifyStatement .
func (sd DeletedAtStampDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		curTime := stmt.DB.NowFunc().UnixMilli()
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

		DeletedatstampQueryclause(sd).ModifyStatement(stmt)
		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.Build(stmt.DB.Callback().Update().Clauses...)
	}
}

// DeletedAtStampCreateClause .
type DeletedAtStampCreateClause struct {
	baseClause

	Field *schema.Field
}

// ModifyStatement .
func (c DeletedAtStampCreateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 && !stmt.Statement.Unscoped {
		stmt.SetColumn(c.Field.Name, sql.NullInt64{Int64: 0, Valid: true})
	}
}
