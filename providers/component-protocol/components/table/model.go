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

package table

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

type (
	// Data .
	Data struct {
		Table      Table                                    `json:"table,omitempty"`
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
	}

	// Table .
	Table struct {
		Columns ColumnsInfo `json:"columns,omitempty"`
		Rows    []Row       `json:"rows,omitempty"`

		PageNo   uint64 `json:"pageNo,omitempty"`
		PageSize uint64 `json:"pageSize,omitempty"`
		Total    uint64 `json:"total,omitempty"`
	}
)

type (
	// ColumnsInfo .
	ColumnsInfo struct {
		// Merges merge some columns into one.
		// +optional
		Merges map[ColumnKey]MergedColumn `json:"merges,omitempty"`

		// Orders is the order of columns.
		// If some columns is merged, just put the columnKey for exhibition.
		Orders []ColumnKey `json:"orders,omitempty"`

		// ColumnsMap contains all columns.
		ColumnsMap map[ColumnKey]Column `json:"columnsMap,omitempty"`
	}

	// MergedColumn .
	MergedColumn struct {
		Orders []ColumnKey `json:"orders,omitempty"`
	}

	// Column .
	Column struct {
		Title string `json:"title,omitempty"`
		Tip   string `json:"tip,omitempty"`

		FieldBindToOrder string `json:"fieldBindToOrder,omitempty"` // bind which field to order
		AscOrder         *bool  `json:"ascOrder,omitempty"`         // true for asc, false for desc, nil for no sort
		EnableSort       bool   `json:"enableSort"`                 // true can sort
		Hidden           bool   `json:"hidden,omitempty"`           // true can hidden
		cptype.Extra
	}

	// ColumnKey .
	ColumnKey string
)

type (
	// Row .
	Row struct {
		ID         RowID              `json:"id,omitempty"` // row id, used for row-level operations
		Selectable bool               `json:"selectable,omitempty"`
		Selected   bool               `json:"selected,omitempty"`
		CellsMap   map[ColumnKey]Cell `json:"cellsMap,omitempty"`

		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
	}

	// RowID .
	RowID string
)

type (
	// Cell .
	Cell struct {
		ID   string          `json:"id,omitempty"`
		Tip  string          `json:"tip,omitempty"`
		Type CellType        `json:"type,omitempty"`
		Data cptype.ExtraMap `json:"data,omitempty"`

		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
		cptype.Extra
	}

	// CellType .
	CellType string
)
