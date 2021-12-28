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

// Icon .
type Icon struct {
	URL  string `json:"url,omitempty"`
	Type string `json:"type,omitempty"`
}

// ModelType .
func (i Icon) ModelType() string { return "icon" }

// NewURLIcon from url, such as: /api/files/{uuid}
func NewURLIcon(url string) *Icon {
	return &Icon{URL: url}
}

// NewTypedIcon from type, such as: ISSUE_ICON.issue.TASK
func NewTypedIcon(_type string) *Icon {
	return &Icon{Type: _type}
}
