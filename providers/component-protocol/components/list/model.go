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

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/components/commodel"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// ItemCommStatus .
type ItemCommStatus string

const (
	// ItemCommStatusDefault default status
	ItemCommStatusDefault ItemCommStatus = ""
	// ItemCommStatusSuccess success status
	ItemCommStatusSuccess ItemCommStatus = "success"
	// ItemCommStatusInfo info status
	ItemCommStatusInfo ItemCommStatus = "info"
	// ItemCommStatusWarning warning status
	ItemCommStatusWarning ItemCommStatus = "warning"
	// ItemCommStatusError error status
	ItemCommStatusError ItemCommStatus = "error"
)

type (
	// Data includes list of items.
	Data struct {
		Title        string                                   `json:"title,omitempty"`
		TitleSummary string                                   `json:"titleSummary,omitempty"`
		List         []Item                                   `json:"list,omitempty"`
		PageNo       uint64                                   `json:"pageNo,omitempty"`
		PageSize     uint64                                   `json:"pageSize,omitempty"`
		Total        uint64                                   `json:"total,omitempty"`
		Operations   map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
		UserIDs      []string                                 `json:"userIDs,omitempty"`
	}

	// Item minimum unit of list
	Item struct {
		// uniq id of the item, e.g: appID, projID, ...
		ID           string `json:"id,omitempty"`
		Title        string `json:"title,omitempty"`
		TitleSummary string `json:"titleSummary,omitempty"`
		// logo link url of title
		//LogoURL          string      `json:"logoURL,omitempty"`
		Star             *bool       `json:"star,omitempty"`
		MainState        *StateInfo  `json:"mainState,omitempty"`
		TitleState       []StateInfo `json:"titleState,omitempty"`
		Description      string      `json:"description,omitempty"`
		BackgroundImgURL string      `json:"backgroundImgURL,omitempty"`
		KvInfos          []KvInfo    `json:"kvInfos,omitempty"`
		Selectable       bool        `json:"selectable"`
		// columns show in the item, e.g user, time
		ColumnsInfo map[string]interface{} `json:"columnsInfo,omitempty"`
		// operations on the frond
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
		// operations folded
		MoreOperations []MoreOpItem   `json:"moreOperations,omitempty"`
		Icon           *commodel.Icon `json:"icon,omitempty"`
		cptype.Extra
	}

	// StateInfo .
	StateInfo struct {
		Text   string         `json:"text,omitempty"`
		Status ItemCommStatus `json:"status,omitempty"`
		// right or left
		SuffixIcon string                                   `json:"suffixIcon,omitempty"`
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations"`
	}

	// KvInfo .
	KvInfo struct {
		ID    string `json:"id,omitempty"`
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
		Icon  string `json:"icon,omitempty"`
		Tip   string `json:"tip,omitempty"`
		// red green etc.
		Color string `json:"color,omitempty"`
		// metaInfo related operations
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations"`
	}

	// MoreOpItem more operation item info
	MoreOpItem struct {
		ID   string `json:"id,omitempty"`
		Text string `json:"text,omitempty"`
		Icon string `json:"icon,omitempty"`
		// more operation related operations
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations"`
	}
)
