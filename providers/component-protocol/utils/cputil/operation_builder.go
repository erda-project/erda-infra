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

package cputil

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// OperationBuilder .
type OperationBuilder struct {
	cptype.Operation
}

// NewOpBuilder .
func NewOpBuilder() *OperationBuilder { return &OperationBuilder{} }

// Build .
func (b *OperationBuilder) Build() cptype.Operation {
	return b.Operation
}

// WithText .
func (b *OperationBuilder) WithText(text string) *OperationBuilder {
	b.Operation.Text = text
	return b
}

// WithConfirmTip .
func (b *OperationBuilder) WithConfirmTip(confirmTip string) *OperationBuilder {
	b.Operation.Confirm = confirmTip
	return b
}

// WithDisable .
func (b *OperationBuilder) WithDisable(disable bool, tip string) *OperationBuilder {
	b.Operation.Disabled = disable
	b.Operation.Tip = tip
	return b
}

// WithTip .
func (b *OperationBuilder) WithTip(tip string) *OperationBuilder {
	b.Operation.Tip = tip
	return b
}

// WithSkipRender .
func (b *OperationBuilder) WithSkipRender(skipRender bool) *OperationBuilder {
	b.Operation.SkipRender = skipRender
	return b
}

// WithAsync .
func (b *OperationBuilder) WithAsync(async bool) *OperationBuilder {
	b.Operation.Async = async
	return b
}

// WithServerDataPtr .
func (b *OperationBuilder) WithServerDataPtr(inputPtr interface{}) *OperationBuilder {
	var serverData cptype.OpServerData
	MustObjJSONTransfer(inputPtr, &serverData)
	b.Operation.ServerData = &serverData
	return b
}
