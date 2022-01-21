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

// ProgressBar is progress bar.
// formula: BarPercent = BarCompletedNum / BarTotalNum * 100
// example: 7/10 * 100 = 70
type ProgressBar struct {
	// BarCompletedNum is Numerator.
	BarCompletedNum int64 `json:"barCompletedNum,omitempty"`
	// BarTotalNum is Denominator.
	BarTotalNum int64 `json:"barTotalNum,omitempty"`

	// BarPercent is the calculated result.
	// Optional.
	// If not present, BarCompletedNum and BarTotalNum must provide.
	// For some situations, backend only have percent value.
	BarPercent float64 `json:"barPercent,omitempty"` // optional, percent, range: [0,100], such as: 0.1, 20, 100

	Text   string        `json:"text,omitempty"`   // optional, bar detail text
	Tip    string        `json:"tip,omitempty"`    // optional, tip
	Status UnifiedStatus `json:"status,omitempty"` // optional, status
}

// ModelType .
func (p ProgressBar) ModelType() string {
	return "progressBar"
}
