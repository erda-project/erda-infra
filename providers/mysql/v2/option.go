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

package v2

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

var (
	// DESC means select by DESC
	DESC Order = "DESC"
	// ASC means select by ASC
	ASC Order = "ASC"
)

// Option is the function processes *gorm.DB
type Option func(db *gorm.DB) *gorm.DB

// Column is an interface for shortcut of WhereColumn, OrderColumn, SetColumn
type Column interface {
	WhereColumn
	OrderColumn
	SetColumn
}

// WhereColumn contains options for conditions
type WhereColumn interface {
	Is(value interface{}) Option
	In(values []interface{}) Option
	InMap(values map[interface{}]struct{}) Option
	Like(value interface{}) Option
	GreaterThan(value interface{}) Option
	EqGreaterThan(value interface{}) Option
	LessThan(value interface{}) Option
	EqLessThan(value interface{}) Option
}

// OrderColumn contains order by options
type OrderColumn interface {
	DESC() Option
	ASC() Option
}

// SetColumn is used in update statements
type SetColumn interface {
	Set(value interface{}) Option
}

// Where is the option for conditions
func Where(format string, args ...interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(format, args...)
	}
}

// Wheres is the option for conditions.
// the m can be a map or a struct.
func Wheres(m interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(m)
	}
}

// Col returns a Column interface
func Col(col string) Column {
	return column{col: col}
}

type column struct {
	col string
}

func (c column) Is(value interface{}) Option {
	if value == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(c.col + " IS NULL")
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(c.col+" = ?", value)
	}
}

func (c column) In(values []interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(c.col+" IN ?", values)
	}
}

func (c column) InMap(values map[interface{}]struct{}) Option {
	var values_ []interface{}
	for key := range values {
		values_ = append(values_, key)
	}
	fmt.Printf("values_: %v", values_)
	return c.In(values_)
}

func (c column) Like(value interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(c.col+" LIKE ?", value)
	}
}

func (c column) GreaterThan(value interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(c.col+" > ?", value)
	}
}

func (c column) EqGreaterThan(value interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(c.col+" >= ?", value)
	}
}

func (c column) LessThan(value interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(c.col+" < ?", value)
	}
}

func (c column) EqLessThan(value interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(c.col+" <= ?", value)
	}
}

func (c column) DESC() Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(c.col + " DESC")
	}
}

func (c column) ASC() Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(c.col + " ASC")
	}
}

func (c column) Set(value interface{}) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Update(c.col, value)
	}
}

// WhereValue contains the option presents where the value in
type WhereValue interface {
	In(cols ...string) Option
}

func Value(value interface{}) WhereValue {
	return whereValue{value: value}
}

type whereValue struct {
	value interface{}
}

func (w whereValue) In(cols ...string) Option {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("? IN (%s)", strings.Join(cols, ",")), w.value)
	}
}

// Paging returns an Option for selecting by paging
func Paging(size, no int) Option {
	if size < 0 {
		size = 0
	}
	if no < 1 {
		no = 1
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(size).Offset((no - 1) * size)
	}
}

// Order .
type Order string

// OrderBy returns an Option for selecting order by some column
func OrderBy(col string, order Order) Option {
	if !strings.EqualFold(string(order), string(DESC)) &&
		!strings.EqualFold(string(order), string(ASC)) {
		order = "DESC"
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(col + " " + strings.ToUpper(string(order)))
	}
}
