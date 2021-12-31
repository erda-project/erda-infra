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

import "sync"

// Below is standard struct for line graph related.
type (
	// Data includes list.
	Data struct {
		Title      string   `json:"title"`
		Dimensions []string `json:"dimensions"`
		XAxis      *Axis    `json:"xAxis"`     // x axis
		YAxis      []*Axis  `json:"yAxis"`     // y axis
		Formatter  string   `json:"formatter"` // data formatter
		Inverse    bool     `json:"inverse"`   // inverted xAxis and yAxis
		sync.RWMutex
	}

	// Axis defined struct.
	Axis struct {
		Dimension string        `json:"dimension,omitempty"` // The xAxis can have no dimensions
		Values    []interface{} `json:"values"`
		Inverse   bool          `json:"inverse"` // inverted values
	}
)

// New .
func New(title string) *Data {
	return &Data{
		Title:      title,
		Dimensions: *new([]string),
		XAxis:      new(Axis),
		YAxis:      *new([]*Axis),
		Inverse:    false,
	}
}

// SetXAxis .
func (d *Data) SetXAxis(values ...interface{}) {
	d.Lock()
	defer d.Unlock()
	d.XAxis.Values = append(d.XAxis.Values, values...)
}

// SetYAxis .
func (d *Data) SetYAxis(dimension string, values ...interface{}) {
	d.Lock()
	defer d.Unlock()
	d.Dimensions = append(d.Dimensions, dimension)
	d.YAxis = append(d.YAxis, &Axis{Dimension: dimension, Values: values})
}
