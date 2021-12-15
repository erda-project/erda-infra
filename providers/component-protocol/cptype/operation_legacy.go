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

package cptype

// LegacyOperation .
type LegacyOperation struct {
	Key         string `json:"key,omitempty"`
	Text        string `json:"text,omitempty"`
	Confirm     string `json:"confirm,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
	DisabledTip string `json:"disabledTip,omitempty"`
	Reload      bool   `json:"reload,omitempty"`
}
