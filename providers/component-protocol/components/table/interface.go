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

// ITable is user-level interface for table.
type ITable interface {
	cptype.IComponent
	ITableStdOps
}

// ITableStdOps defines all standard operations of table.
type ITableStdOps interface {
	// table-level
	RegisterTableChangePageOp(opData OpTableChangePage) (opFunc cptype.OperationFunc)
	RegisterTableSortOp(opData OpTableChangeSort) (opFunc cptype.OperationFunc)
	RegisterBatchRowsHandleOp(opData OpBatchRowsHandle) (opFunc cptype.OperationFunc)
	// row-level
	RegisterRowSelectOp(opData OpRowSelect) (opFunc cptype.OperationFunc)
	RegisterRowAddOp(opData OpRowAdd) (opFunc cptype.OperationFunc)
	RegisterRowEditOp(opData OpRowEdit) (opFunc cptype.OperationFunc)
	RegisterRowDeleteOp(opData OpRowDelete) (opFunc cptype.OperationFunc)
	// cell-level
}
