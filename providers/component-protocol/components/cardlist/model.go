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

package cardlist

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

type (
	// Data card std data
	Data struct {
		Total        int    `json:"total"`
		Title        string `json:"title"`
		TitleSummary string `json:"titleSummary"`
		Cards        []Card `json:"cards,omitempty"`
	}

	// Card .
	Card struct {
		ID         string       `json:"id"`
		ImgURL     string       `json:"imgURL"`
		Icon       string       `json:"icon"`
		Title      string       `json:"title"`
		TitleState []TitleState `json:"titleState"`
		Star       bool         `json:"star"`
		TextMeta   []TextMeta   `json:"textMeta"`
		cptype.Extra
	}

	// TitleState .
	TitleState struct {
		Text   string `json:"text"`
		Color  string `json:"color"`
		Status string `json:"status"`
	}

	// TextMeta .
	TextMeta struct {
		MainText   float64                                  `json:"mainText"`
		SubText    string                                   `json:"subText"`
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations"`
	}
)
