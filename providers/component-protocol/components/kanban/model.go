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

// Below is standard struct for kanban related.
type (
	// Data includes list of boards.
	Data struct {
		Boards     []Board                                  `json:"boards,omitempty"`
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
		UserIDs    []string                                 `json:"userIDs,omitempty"`
	}
	// Board includes list of cards.
	Board struct {
		ID         string                                   `json:"id,omitempty"`
		Title      string                                   `json:"title,omitempty"`
		Cards      []Card                                   `json:"cards,omitempty"`
		PageNo     uint64                                   `json:"pageNo,omitempty"`
		PageSize   uint64                                   `json:"pageSize,omitempty"`
		Total      uint64                                   `json:"total,omitempty"`
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
	}
	// Card is a minimal unit in kanban.
	Card struct {
		ID         string                                   `json:"id,omitempty"`
		Title      string                                   `json:"title,omitempty"`
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`

		cptype.Extra
	}
)
