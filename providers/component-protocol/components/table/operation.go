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
)

// table-level
type (
	// OpTableChangePage includes std op and server&client data.
	OpTableChangePage struct {
		cptype.Operation
		ServerData OpTableChangePageServerData `json:"serverData,omitempty"`
		ClientData OpTableChangePageClientData `json:"clientData,omitempty"`
	}
	// OpTableChangePageServerData .
	OpTableChangePageServerData struct {
	}
	// OpTableChangePageClientData .
	OpTableChangePageClientData struct {
		Title string `json:"title,omitempty"`
	}
)

// OpKey .
func (o OpTableChangePage) OpKey() cptype.OperationKey { return "changePage" }

type (
	// OpTableChangeSort .
	OpTableChangeSort struct {
		cptype.Operation
		ServerData OpTableChangePageServerData `json:"serverData,omitempty"`
		ClientData OpTableChangeSortClientData `json:"clientData,omitempty"`
	}
	// OpTableChangeSortServerData .
	OpTableChangeSortServerData struct {
	}
	// OpTableChangeSortClientData .
	OpTableChangeSortClientData struct {
		DataRef *Column `json:"dataRef,omitempty"`
	}
)

// OpKey .
func (o OpTableChangeSort) OpKey() cptype.OperationKey {
	return "changeSort"
}

type (
	// OpBatchRowsHandle .
	OpBatchRowsHandle struct {
		cptype.Operation
		ServerData OpTableChangePageServerData `json:"serverData,omitempty"`
		ClientData OpBatchRowsHandleClientData `json:"clientData,omitempty"`
	}
	// OpBatchRowsHandleServerData .
	OpBatchRowsHandleServerData struct {
		Options []OpBatchRowsHandleOption `json:"options,omitempty"`
	}
	// OpBatchRowsHandleOption .
	OpBatchRowsHandleOption struct {
		ID              string         `json:"id,omitempty"`
		Text            string         `json:"text,omitempty"`
		Icon            *commodel.Icon `json:"icon,omitempty"`
		AllowedRowIDs   []string       `json:"allowedRowIDs,omitempty"`
		ForbiddenRowIDs []string       `json:"forbiddenRowIDs,omitempty"`
	}
	// OpBatchRowsHandleClientData .
	OpBatchRowsHandleClientData struct {
		DataRef          *OpBatchRowsHandleOption `json:"dataRef,omitempty"`
		SelectedOptionID string                   `json:"selectedOptionID,omitempty"`
		SelectedRowIDs   []string                 `json:"selectedRowIDs,omitempty"`
	}
)

// OpKey .
func (o OpBatchRowsHandle) OpKey() cptype.OperationKey {
	return "batchRowsHandle"
}

// row-level
type (
	// OpRowSelect .
	OpRowSelect struct {
		cptype.Operation
		ServerData OpRowSelectServerData `json:"serverData,omitempty"`
		ClientData OpRowSelectClientData `json:"clientData,omitempty"`
	}
	// OpRowSelectServerData .
	OpRowSelectServerData struct {
	}
	// OpRowSelectClientData .
	OpRowSelectClientData struct {
		DataRef *Row `json:"dataRef,omitempty"`
	}
)

// OpKey .
func (o OpRowSelect) OpKey() cptype.OperationKey {
	return "rowSelect"
}

type (
	// OpRowAdd .
	OpRowAdd struct {
		cptype.Operation
		ServerData OpRowAddServerData `json:"serverData,omitempty"`
		ClientData OpRowAddClientData `json:"clientData,omitempty"`
	}
	// OpRowAddServerData .
	OpRowAddServerData struct {
	}
	// OpRowAddClientData .
	OpRowAddClientData struct {
		LastRowID string `json:"lastRowID,omitempty"`
		NextRowID string `json:"nextRowID,omitempty"`
		NewRow    *Row   `json:"newRow,omitempty"`
	}
)

// OpKey .
func (o OpRowAdd) OpKey() cptype.OperationKey {
	return "rowAdd"
}

type (
	// OpRowEdit .
	OpRowEdit struct {
		cptype.Operation
		ServerData OpRowEditServerData `json:"serverData,omitempty"`
		ClientData OpRowEditClientData `json:"clientData,omitempty"`
	}
	// OpRowEditServerData .
	OpRowEditServerData struct {
	}
	// OpRowEditClientData .
	OpRowEditClientData struct {
		DataRef *Row `json:"dataRef,omitempty"`
		NewRow  *Row `json:"newRow,omitempty"`
	}
)

// OpKey .
func (o OpRowEdit) OpKey() cptype.OperationKey {
	return "rowEdit"
}

type (
	// OpRowDelete .
	OpRowDelete struct {
		cptype.Operation
		ServerData OpRowDeleteServerData `json:"serverData,omitempty"`
		ClientData OpRowDeleteClientData `json:"clientData,omitempty"`
	}
	// OpRowDeleteServerData .
	OpRowDeleteServerData struct {
	}
	// OpRowDeleteClientData .
	OpRowDeleteClientData struct {
		DataRef *Row `json:"dataRef,omitempty"`
	}
)

// OpKey .
func (o OpRowDelete) OpKey() cptype.OperationKey {
	return "rowDelete"
}

// cell-level
