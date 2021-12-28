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

// DropDownMenu .
type DropDownMenu struct {
	Menus []MenuItem `json:"menus,omitempty"`

	Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
}

// MenuItem .
type MenuItem struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`

	Icon          Icon          `json:"icon,omitempty"`
	UnifiedStatus UnifiedStatus `json:"unifiedStatus,omitempty"`

	Selected bool `json:"selected,omitempty"`
	Disabled bool `json:"disabled,omitempty"`
	Hidden   bool `json:"hidden,omitempty"`
	Tip      bool `json:"tip,omitempty"`
}

type (
	// OpDropDownMenuChange .
	OpDropDownMenuChange struct {
		cptype.Operation
		ServerData OpDropDownMenuChangeServerData `json:"serverData,omitempty"`
		ClientData OpDropDownMenuChangeClientData `json:"clientData,omitempty"`
	}
	// OpDropDownMenuChangeServerData .
	OpDropDownMenuChangeServerData struct{}
	// OpDropDownMenuChangeClientData .
	OpDropDownMenuChangeClientData struct {
		DataRef        *MenuItem `json:"dataRef,omitempty"`
		SelectedItemID string    `json:"selectedItemID,omitempty"`
	}
)
