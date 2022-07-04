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
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// ErrInvalidTransaction means the transaction is alread committed or roll backed
var ErrInvalidTransaction = errors.New("invalid transaction, it is already committed or roll backed")

// TX contains the CRUS APIs
type TX struct {
	Error error

	db    *gorm.DB
	inTx  bool
	valid bool
}

// NewTx returns a *TX
func NewTx(db *gorm.DB) *TX {
	return &TX{db: db, valid: true}
}

// Create inserts a row
func (tx *TX) Create(i interface{}) error {
	if tx.inTx && !tx.valid {
		return ErrInvalidTransaction
	}
	tx.Error = tx.db.Create(i).Error
	return tx.Error
}

// CreateInBatches inserts multi rows
func (tx *TX) CreateInBatches(i interface{}, size int) error {
	if tx.inTx && !tx.valid {
		return ErrInvalidTransaction
	}
	tx.Error = tx.db.CreateInBatches(i, size).Error
	return tx.Error
}

// Delete deletes rows with conditions
func (tx *TX) Delete(i interface{}, options ...Option) (int64, error) {
	if tx.inTx && !tx.valid {
		return 0, ErrInvalidTransaction
	}
	var db = tx.DB()
	for _, opt := range options {
		db = opt(db)
	}
	db = db.Delete(i)
	return db.RowsAffected, db.Error
}

// Updates updates the model i with the given value. v can be a map or a model struct.
// options is conditions.
func (tx *TX) Updates(i, v interface{}, options ...Option) error {
	if tx.inTx && !tx.valid {
		return ErrInvalidTransaction
	}
	var db = tx.DB()
	for _, opt := range options {
		db = opt(db)
	}
	return db.Model(i).Updates(v).Error
}

// SetColumns is used to set columns.
// At least one SetColumn Option in the options.
func (tx *TX) SetColumns(i interface{}, options ...Option) error {
	if tx.inTx && !tx.valid {
		return ErrInvalidTransaction
	}
	var db = tx.DB()
	db = db.Model(i)
	for _, opt := range options {
		db = opt(db)
	}
	return db.Error
}

// List lists records.
func (tx *TX) List(i interface{}, options ...Option) (int64, error) {
	var total int64
	var db = tx.DB()
	for _, opt := range options {
		db = opt(db)
	}

	err := db.Find(i).Count(&total).Error
	if err == nil {
		return total, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return 0, err
}

// Get gets the first record.
func (tx *TX) Get(i interface{}, options ...Option) (bool, error) {
	var db = tx.DB()
	for _, opt := range options {
		db = opt(db)
	}

	err := db.First(i).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// Commit commits the transaction.
func (tx *TX) Commit() error {
	if !tx.inTx {
		return errors.New("not in transaction")
	}
	if !tx.valid {
		return ErrInvalidTransaction
	}
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "can not commit with error")
	}
	tx.db.Commit()
	tx.valid = false
	return nil
}

// Rollback rollbacks the transaction.
func (tx *TX) Rollback() error {
	if !tx.inTx {
		return errors.New("not in transaction")
	}
	if !tx.valid {
		return ErrInvalidTransaction
	}
	tx.db.Rollback()
	tx.valid = false
	return nil
}

// CommitOrRollback commits the transaction if db.Error is nil,
// or rollbacks if the db.Error is not nil.
func (tx *TX) CommitOrRollback() {
	if tx.inTx && !tx.valid {
		return
	}
	if tx.Error == nil {
		tx.db.Commit()
	} else {
		tx.db.Rollback()
	}
	tx.valid = false
}

// DB returns the raw *gorm.DB
func (tx *TX) DB() *gorm.DB {
	return tx.db
}
