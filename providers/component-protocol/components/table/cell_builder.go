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
	"github.com/erda-project/erda-infra/providers/component-protocol/components/commodel"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

type (
	// CellBuilder .
	CellBuilder struct {
		*Cell
	}
	// ITypedCellBuilder used to not export typedCellBuilder
	ITypedCellBuilder interface {
		getCellBuilder() *CellBuilder
		Build() Cell
	}
	// typedCellBuilder .
	typedCellBuilder struct {
		*CellBuilder
	}
)

func (tb *typedCellBuilder) getCellBuilder() *CellBuilder {
	return tb.CellBuilder
}

func newCellBuilder() *CellBuilder {
	return &CellBuilder{Cell: &Cell{}}
}

func (cb *CellBuilder) typed() *typedCellBuilder {
	return &typedCellBuilder{CellBuilder: cb}
}

// Build .
func (tb *typedCellBuilder) Build() Cell {
	return *tb.Cell
}

// WithID .
func (tb *typedCellBuilder) WithID(id string) *typedCellBuilder {
	tb.Cell.ID = id
	return tb
}

// WithTip .
func (tb *typedCellBuilder) WithTip(tip string) *typedCellBuilder {
	tb.Cell.Tip = tip
	return tb
}

// WithOperations .
func (tb *typedCellBuilder) WithOperations(ops map[cptype.OperationKey]cptype.Operation) *typedCellBuilder {
	tb.Cell.Operations = ops
	return tb
}

// WithOperation .
func (tb *typedCellBuilder) WithOperation(opKey cptype.OperationKey, op cptype.Operation) *typedCellBuilder {
	if tb.Cell.Operations == nil {
		tb.Cell.Operations = make(map[cptype.OperationKey]cptype.Operation)
	}
	tb.Cell.Operations[opKey] = op
	return tb
}

// NewTextCell .
func NewTextCell(text string) ITypedCellBuilder {
	cb := newCellBuilder()
	cb.Cell.Type = CellType(commodel.Text{}.ModelType())
	cputil.MustObjJSONTransfer(&commodel.Text{Text: text}, &cb.Data)
	return cb.typed()
}

// NewKVCell .
func NewKVCell(k, v string) ITypedCellBuilder {
	cb := newCellBuilder()
	cb.Cell.Type = CellType(commodel.KV{}.ModelType())
	cputil.MustObjJSONTransfer(&commodel.KV{K: k, V: v}, &cb.Data)
	return cb.typed()
}

// NewIconCell .
func NewIconCell(icon commodel.Icon) ITypedCellBuilder {
	cb := newCellBuilder()
	cb.Cell.Type = CellType(commodel.Icon{}.ModelType())
	cputil.MustObjJSONTransfer(commodel.NewTypedIcon("ISSUE_ICON.issue.TASK"), &cb.Data)
	return cb.typed()
}

// NewUserCell .
func NewUserCell(user commodel.User) ITypedCellBuilder {
	cb := newCellBuilder()
	cb.Cell.Type = CellType(commodel.User{}.ModelType())
	cputil.MustObjJSONTransfer(&commodel.User{ID: "1", Name: "Bob"}, &cb.Data)
	return cb.typed()
}

// NewLabelsCell .
func NewLabelsCell(labels commodel.Labels) ITypedCellBuilder {
	cb := newCellBuilder()
	cb.Cell.Type = CellType(commodel.Labels{}.ModelType())
	cputil.MustObjJSONTransfer(&labels, &cb.Data)
	return cb.typed()
}

// NewDropDownMenuCell .
func NewDropDownMenuCell(menu commodel.DropDownMenu) ITypedCellBuilder {
	cb := newCellBuilder()
	cb.Cell.Type = CellType(commodel.DropDownMenu{}.ModelType())
	cputil.MustObjJSONTransfer(&menu, &cb.Data)
	return cb.typed()
}

// NewUserSelectorCell .
func NewUserSelectorCell(selector commodel.UserSelector) ITypedCellBuilder {
	cb := newCellBuilder()
	cb.Cell.Type = CellType(commodel.UserSelector{}.ModelType())
	cputil.MustObjJSONTransfer(&selector, &cb.Data)
	return cb.typed()
}
