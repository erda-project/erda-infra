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

package list

import "github.com/erda-project/erda-infra/providers/component-protocol/cptype"

const (
	// OpListChangePageKey list level
	OpListChangePageKey = "changePage"

	// OpItemStarKey item level
	OpItemStarKey = "star"
	// OpItemClickGotoKey item level
	OpItemClickGotoKey = "clickGoto"
	// OpItemClickKey .
	OpItemClickKey = "click"
)

// list-level
type (
	// OpChangePage paging operation data
	OpChangePage struct {
		cptype.Operation
		ServerData OpChangePageServerData `json:"serverData,omitempty"`
		ClientData OpChangePageClientData `json:"clientData,omitempty"`
	}

	// OpChangePageServerData server data
	OpChangePageServerData struct{}

	// OpChangePageClientData data
	OpChangePageClientData struct {
		PageNo   uint64 `json:"pageNo,omitempty"`
		PageSize uint64 `json:"pageSize,omitempty"`
	}
)

// OpKey .
func (o OpChangePage) OpKey() cptype.OperationKey { return OpListChangePageKey }

type (
	// OpItemStar .
	OpItemStar struct {
		cptype.Operation
		ServerData OpItemStarServerData `json:"serverData,omitempty"`
		ClientData OpItemStarClientData `json:"clientData,omitempty"`
	}

	// OpItemStarServerData server data
	OpItemStarServerData struct{}

	// OpItemStarClientData data
	OpItemStarClientData struct {
		DataRef *Item `json:"dataRef,omitempty"`
	}
)

// OpKey .
func (o OpItemStar) OpKey() cptype.OperationKey { return OpItemStarKey }

type (
	// OpItemClickGoto .
	OpItemClickGoto struct {
		cptype.Operation
		ServerData OpItemClickGotoServerData `json:"serverData,omitempty"`
		ClientData OpItemClickGotoClientData `json:"clientData,omitempty"`
	}

	// OpItemClickGotoServerData server data
	OpItemClickGotoServerData struct {
		OpItemBasicServerData
	}

	// OpItemClickGotoClientData data
	OpItemClickGotoClientData struct{}
)

// OpItemBasicServerData .
type OpItemBasicServerData struct {
	// open a new fronted page
	JumpOut bool `json:"jumpOut,omitempty"`
	// the jump out target of the new opened fronted page, e.g: projectAllIssue page
	Target string `json:"target,omitempty"`
	// params needed for jumping to the new page, e.g: projectId
	Params map[string]interface{} `json:"params,omitempty"`
	// the query condition of the new page, e.g: issueFilter__urlQuery
	Query map[string]interface{} `json:"query,omitempty"`
}

// OpKey .
func (o OpItemClickGoto) OpKey() cptype.OperationKey { return OpItemClickGotoKey }

type (
	// OpItemClick .
	OpItemClick struct {
		cptype.Operation
		ServerData OpItemClickServerData `json:"serverData,omitempty"`
		ClientData OpItemClickClientData `json:"clientData,omitempty"`
	}

	// OpItemClickServerData server data
	OpItemClickServerData struct {
		OpItemBasicServerData
	}

	// OpItemClickClientData data
	OpItemClickClientData struct {
		DataRef *Item `json:"dataRef,omitempty"`
	}
)

// OpKey .
func (o OpItemClick) OpKey() cptype.OperationKey { return OpItemClickKey }
