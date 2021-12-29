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

// Labels .
type Labels struct {
	Labels []Label `json:"labels,omitempty"`

	Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
}

// ModelType .
func (l Labels) ModelType() string { return "labels" }

// Label .
type Label struct {
	ID    string       `json:"id,omitempty"`
	Title string       `json:"title,omitempty"`
	Color UnifiedColor `json:"color,omitempty"`

	Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
}

// ModelType .
func (l Label) ModelType() string { return "label" }
