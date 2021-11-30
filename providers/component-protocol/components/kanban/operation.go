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

package kanban

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// kanban-level
type (
	// OpBoardCreate includes std op and server&client data.
	OpBoardCreate struct {
		cptype.Operation
		ServerData OpBoardCreateServerData `json:"serverData,omitempty"`
		ClientData OpBoardCreateClientData `json:"clientData,omitempty"`
	}
	// OpBoardCreateServerData .
	OpBoardCreateServerData struct{}
	// OpBoardCreateClientData .
	OpBoardCreateClientData struct {
		Title string `json:"title,omitempty"`
	}
)

// OpKey .
func (o OpBoardCreate) OpKey() cptype.OperationKey { return "boardCreate" }

// board-level
type (
	// OpBoardLoadMore .
	OpBoardLoadMore struct {
		cptype.Operation
		ServerData OpBoardLoadMoreServerData `json:"serverData,omitempty"`
		ClientData OpBoardLoadMoreClientData `json:"clientData,omitempty"`
	}
	// OpBoardLoadMoreServerData .
	OpBoardLoadMoreServerData struct{}
	// OpBoardLoadMoreClientData .
	OpBoardLoadMoreClientData struct {
		DataRef *Board `json:"dataRef,omitempty"`

		PageNo   uint64 `json:"pageNo,omitempty"`
		PageSize uint64 `json:"pageSize,omitempty"`
	}
)

// OpKey .
func (o OpBoardLoadMore) OpKey() cptype.OperationKey { return "boardLoadMore" }

type (
	// OpBoardUpdate .
	OpBoardUpdate struct {
		cptype.Operation
		ServerData OpBoardUpdateServerData `json:"serverData,omitempty"`
		ClientData OpBoardUpdateClientData `json:"clientData,omitempty"`
	}
	// OpBoardUpdateServerData .
	OpBoardUpdateServerData struct{}
	// OpBoardUpdateClientData .
	OpBoardUpdateClientData struct {
		DataRef *Board `json:"dataRef,omitempty"`

		// update fields
		Title string `json:"title,omitempty"`
	}
)

// OpKey .
func (o OpBoardUpdate) OpKey() cptype.OperationKey { return "boardUpdate" }

type (
	// OpBoardDelete .
	OpBoardDelete struct {
		cptype.Operation
		ServerData OpBoardDeleteServerData `json:"serverData,omitempty"`
		ClientData OpBoardDeleteClientData `json:"clientData,omitempty"`
	}
	// OpBoardDeleteServerData .
	OpBoardDeleteServerData struct{}
	// OpBoardDeleteClientData .
	OpBoardDeleteClientData struct {
		DataRef *Board `json:"dataRef,omitempty"`
	}
)

// OpKey .
func (o OpBoardDelete) OpKey() cptype.OperationKey { return "boardDelete" }

// card-level
type (
	// OpCardMoveTo .
	OpCardMoveTo struct {
		cptype.Operation
		ServerData OpCardMoveToServerData `json:"serverData,omitempty"`
		ClientData OpCardMoveToClientData `json:"clientData,omitempty"`
	}
	// OpCardMoveToServerData .
	OpCardMoveToServerData struct {
		AllowedTargetBoardIDs []string `json:"allowedTargetBoardIDs,omitempty"`
	}
	// OpCardMoveToClientData .
	OpCardMoveToClientData struct {
		DataRef *Card `json:"dataRef,omitempty"`

		TargetBoardID string `json:"targetBoardID,omitempty"`
	}
)

// OpKey .
func (o OpCardMoveTo) OpKey() cptype.OperationKey { return "cardMoveTo" }
