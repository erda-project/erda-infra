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

package filter

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

type (
	OpFilter struct {
		cptype.Operation
		ClientData OpFilterClientData `json:"clientData,omitempty"`
	}
	OpFilterClientData struct {
		Values cptype.ExtraMap `json:"values,omitempty"`
	}
)

func (o OpFilter) OpKey() cptype.OperationKey { return "filter" }

type (
	OpFilterItemSave struct {
		cptype.Operation
		ClientData OpFilterItemSaveClientData `json:"clientData,omitempty"`
	}
	OpFilterItemSaveClientData struct {
		Values cptype.ExtraMap `json:"values,omitempty"`
		Label  string          `json:"label,omitempty"`
	}
)

func (o OpFilterItemSave) OpKey() cptype.OperationKey { return "saveFilterSet" }

type (
	OpFilterItemDelete struct {
		cptype.Operation
		ClientData OpFilterItemDeleteClientData `json:"clientData,omitempty"`
	}
	OpFilterItemDeleteClientData struct {
		DataRef *FilterItem `json:"dataRef,omitempty"`
	}
)

func (o OpFilterItemDelete) OpKey() cptype.OperationKey { return "deleteFilterSet" }
