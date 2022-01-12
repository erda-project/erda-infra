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

package commodel

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// MoreOperations .
type MoreOperations struct {
	Ops []MoreOpItem `json:"ops,omitempty"` // use list for order
}

// MoreOpItem more operation item info
type MoreOpItem struct {
	ID   string `json:"id,omitempty"`
	Text string `json:"text,omitempty"`
	Icon *Icon  `json:"icon,omitempty"`

	Operations map[cptype.OperationKey]cptype.Operation `json:"operations"`
}

// ModelType .
func (m MoreOperations) ModelType() string { return "moreOperations" }

type (
	// OpMoreOperationsItemClick .
	OpMoreOperationsItemClick struct {
		OpClick
		ServerData OpMoreOperationsItemClickServerData `json:"serverData,omitempty"`
		ClientData OpMoreOperationsItemClickClientData `json:"clientData,omitempty"`
	}
	// OpMoreOperationsItemClickServerData .
	OpMoreOperationsItemClickServerData struct {
	}
	// OpMoreOperationsItemClickClientData .
	OpMoreOperationsItemClickClientData struct {
		DataRef       *MoreOpItem `json:"dataRef,omitempty"`
		ParentDataRef interface{} `json:"parentDataRef,omitempty"` // optional, such like list row data, table row data
	}
)

type (
	// OpMoreOperationsItemClickGoto .
	OpMoreOperationsItemClickGoto struct {
		OpClickGoto
	}
	// OpMoreOperationsItemClickGotoServerData .
	OpMoreOperationsItemClickGotoServerData struct {
		OpClickGotoServerData
	}
	// OpMoreOperationsItemClickGotoClientData .
	OpMoreOperationsItemClickGotoClientData struct {
		OpMoreOperationsItemClickClientData
	}
)
