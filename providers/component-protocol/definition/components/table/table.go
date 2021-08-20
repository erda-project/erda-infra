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
	"github.com/erda-project/erda-infra/providers/component-protocol/definition/cptype"
)

// CommonTable .
type CommonTable struct {
	State      State                                    `json:"state"`
	Type       string                                   `json:"type"`
	Operations map[cptype.OperationKey]cptype.Operation `json:"operations"`
	Props      Props                                    `json:"props"`
}

// Operate .
type Operate struct {
	RenderType string                `json:"renderType"`
	Value      string                `json:"value"`
	Operations map[string]Operations `json:"operations"`
}

// Operations .
type Operations struct {
	Key     string `json:"key"`
	Text    string `json:"text"`
	Reload  bool   `json:"reload"`
	Confirm string `json:"confirm"`
	Meta    Meta   `json:"meta"`
}

// Meta .
type Meta struct {
	ID int64 `json:"id"`
}

// FormData .
type FormData struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Items    []string `json:"items"`
	Target   string   `json:"target"`
	Channels []string `json:"channels"`
}

// State .
type State struct {
	EditID    uint64 `json:"editId"`
	Operation string `json:"operation"`
	Visible   bool   `json:"visible"`
}

// Props .
type Props struct {
	RowKey     string       `json:"rowKey"`
	Columns    []PropColumn `json:"columns"`
	Pagination bool         `json:"pagination"`
}

// PropColumn .
type PropColumn struct {
	Title     string `json:"title"`
	DataIndex string `json:"dataIndex"`
	Width     uint64 `json:"width"`
}
