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

// UserSelector .
type UserSelector struct {
	Scope           string   `json:"scope,omitempty"` // org/project/app/etc
	SelectedUserIDs []string `json:"selectedUserIDs,omitempty"`

	Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
}

// ModelType .
func (us UserSelector) ModelType() string { return "userSelector" }

type (
	// OpUserSelectorChange .
	OpUserSelectorChange struct {
		cptype.Operation
		ServerData OpUserSelectorChangeServerData `json:"serverData,omitempty"`
		ClientData OpUserSelectorChangeClientData `json:"clientData,omitempty"`
	}
	// OpUserSelectorChangeServerData .
	OpUserSelectorChangeServerData struct{}
	// OpUserSelectorChangeClientData .
	OpUserSelectorChangeClientData struct {
		DataRef         *MenuItem `json:"dataRef,omitempty"`
		SelectedUserIDs []string  `json:"selectedUserIDs,omitempty"`
	}
)
