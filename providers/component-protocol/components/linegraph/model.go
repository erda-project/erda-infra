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

package linegraph

// Below is standard struct for line graph related.
type (
	// Data includes list.
	Data struct {
		Title      string   `json:"title"`
		Dimensions []string `json:"dimensions"`
		XAxis      []*Axis  `json:"xAxis"`   // x axis
		YAxis      []*Axis  `json:"yAxis"`   // y axis
		Inverse    bool     `json:"inverse"` // inverted xAxis and yAxis
	}

	// Axis defined struct.
	Axis struct {
		Dimension string        `json:"dimension,omitempty"`
		Values    []interface{} `json:"values"`
		Inverse   bool          `json:"inverse"` // inverted values
	}
)
