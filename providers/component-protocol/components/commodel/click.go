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

// OpClick
type (
	// OpClick is a simple click op.
	OpClick struct{}
)

// OpKey .
func (o OpClick) OpKey() cptype.OperationKey { return "click" }

// OpClickGoto
type (
	// OpClickGoto means click and goto target link.
	OpClickGoto struct {
	}
	// OpClickGotoServerData .
	OpClickGotoServerData struct {
		// open a new fronted page
		JumpOut bool `json:"jumpOut,omitempty"`
		// the jump out target of the new opened fronted page, e.g: projectAllIssue page
		Target string `json:"target,omitempty"`
		// params needed for jumping to the new page, e.g: projectId
		Params map[string]interface{} `json:"params,omitempty"`
		// the query condition of the new page, e.g: issueFilter__urlQuery
		Query map[string]interface{} `json:"query,omitempty"`
	}
)

// OpKey .
func (o OpClickGoto) OpKey() cptype.OperationKey { return "clickGoto" }
