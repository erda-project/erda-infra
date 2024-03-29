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

// Duration A Duration represents the elapsed time between two instants
// Such as 10s 360s, if the value is -1 means not started
type Duration struct {
	Value int64  `json:"value,omitempty"`
	Tip   string `json:"tip,omitempty"` // Hover on the text, use '\n' for newlines
}

// ModelType .
func (d Duration) ModelType() string { return "duration" }
