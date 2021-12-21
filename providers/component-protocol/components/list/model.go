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

type ItemLabelStatus string

const (
	ItemLabelStatusDefault ItemLabelStatus = ""
	ItemLabelStatusSuccess ItemLabelStatus = "success"
	ItemLabelStatusInfo    ItemLabelStatus = "info"
	ItemLabelStatusWarning ItemLabelStatus = "warning"
	ItemLabelStatusError   ItemLabelStatus = "error"
)

type (
	// Data includes list of items.
	Data struct {
		List       []Item                                   `json:"list,omitempty"`
		PageNo     uint64                                   `json:"pageNo,omitempty"`
		PageSize   uint64                                   `json:"pageSize,omitempty"`
		Total      uint64                                   `json:"total,omitempty"`
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
		UserIDs    []string                                 `json:"userIDs,omitempty"`
	}

	// Item minimum unit of list
	Item struct {
		// uniq id of the item, e.g: appID, projID, ...
		ID    string `json:"id,omitempty"`
		Title string `json:"title,omitempty"`
		// logo link url of title
		LogoURL          string      `json:"LogoURL,omitempty"`
		Star             bool        `json:"star,omitempty"`
		Labels           []ItemLabel `json:"labels,omitempty"`
		Description      string      `json:"description,omitempty"`
		BackgroundImgURL string      `json:"backgroundImgURL,omitempty"`
		MetaInfos        []MetaInfo  `json:"metaInfos,omitempty"`
		// operations on the frond
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
		// operations folded
		MoreOperations MoreOperations `json:"moreOperations,omitempty"`
		cptype.Extra
	}

	// ItemLabel .
	ItemLabel struct {
		Label string `json:"label,omitempty"`
		// optional: label color
		Color string `json:"color,omitempty"`
		// optional: default[gray], success[green], info[blue], warning[yellow], error[red]
		Status ItemLabelStatus `json:"status,omitempty"`
	}

	// MetaInfo .
	MetaInfo struct {
		Label string `json:"label,omitempty"`
		Value string `json:"value,omitempty"`
		Icon  string `json:"icon,omitempty"`
		Tip   string `json:"tip,omitempty"`
		// metaInfo related operations
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations"`
	}

	// MoreOperations .
	MoreOperations struct {
		Operations      map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
		OperationsOrder []cptype.OperationKey                    `json:"operationsOrder,omitempty"`
	}
)
