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

package topn

// Below is standard struct for top related.
type (
	// Data includes list.
	Data struct {
		List []Record `json:"list,omitempty"`
	}

	// Item .
	Item struct {
		ID    string  `json:"id,omitempty"`
		Name  string  `json:"name,omitempty"`
		Value float64 `json:"value,omitempty"`
		Total float64 `json:"total,omitempty"`
		Unit  string  `json:"unit,omitempty"`
	}

	// Record includes one record of top
	Record struct {
		Title          string `json:"title,omitempty"`          // record title
		Items          []Item `json:"items,omitempty"`          // records
		TitleIcon      string `json:"titleIcon,omitempty"`      // title icon
		BackgroundIcon string `json:"backgroundIcon,omitempty"` // record icon
		Span           string `json:"span,omitempty"`           // scale of each record
	}
)
