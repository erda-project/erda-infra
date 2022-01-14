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

import (
	"sync"

	structure "github.com/erda-project/erda-infra/providers/component-protocol/components/commodel/data-structure"
)

// Below is standard struct for line graph related.
type (
	// Data includes list.
	Data struct {
		Title      string     `json:"title"`
		SubTitle   string     `json:"subTitle"`
		Dimensions []string   `json:"dimensions"`
		XAxis      *Axis      `json:"xAxis"` // x axis
		YAxis      []*Axis    `json:"yAxis"` // y axis
		XOptions   *Options   `json:"xOptions"`
		YOptions   []*Options `json:"yOptions"`
		Inverse    bool       `json:"inverse"` // inverted xAxis and yAxis
		sync.RWMutex
	}

	// Axis defined struct.
	Axis struct {
		Dimension string        `json:"dimension,omitempty"` // The xAxis can have no dimensions
		Values    []interface{} `json:"values"`
	}

	// Options .
	Options struct {
		Dimension string                   `json:"dimension,omitempty"`
		Structure *structure.DataStructure `json:"structure"`
		Inverse   bool                     `json:"inverse"` // inverted values
	}

	// DataBuilder .
	DataBuilder struct {
		data *Data
	}

	// OptionsBuilder .
	OptionsBuilder struct {
		options *Options
	}
)

// NewOptionsBuilder .
func NewOptionsBuilder() *OptionsBuilder {
	return &OptionsBuilder{options: &Options{Structure: &structure.DataStructure{}}}
}

// WithDimension .
func (o *OptionsBuilder) WithDimension(dimension string) *OptionsBuilder {
	o.options.Dimension = dimension
	return o
}

// WithType .
func (o *OptionsBuilder) WithType(dataType structure.Type) *OptionsBuilder {
	o.options.Structure.Type = dataType
	return o
}

// WithPrecision .
func (o *OptionsBuilder) WithPrecision(precision structure.Precision) *OptionsBuilder {
	o.options.Structure.Precision = precision
	return o
}

// Build .
func (o *OptionsBuilder) Build() *Options {
	return o.options
}

// NewDataBuilder .
func NewDataBuilder() *DataBuilder {
	return &DataBuilder{data: &Data{Dimensions: []string{}, XAxis: &Axis{}, YAxis: []*Axis{}, XOptions: &Options{}, YOptions: []*Options{}}}
}

// WithTitle .
func (d *DataBuilder) WithTitle(title string) *DataBuilder {
	d.data.Title = title
	return d
}

// WithXAxis .
func (d *DataBuilder) WithXAxis(values ...interface{}) *DataBuilder {
	d.data.XAxis.Values = append(d.data.XAxis.Values, values...)
	return d
}

// WithYAxis .
func (d *DataBuilder) WithYAxis(dimension string, values ...interface{}) *DataBuilder {
	d.data.Dimensions = append(d.data.Dimensions, dimension)
	d.data.YAxis = append(d.data.YAxis, &Axis{Dimension: dimension, Values: values})
	return d
}

// WithXOptions .
func (d *DataBuilder) WithXOptions(options *Options) *DataBuilder {
	d.data.XOptions = options
	return d
}

// WithYOptions .
func (d *DataBuilder) WithYOptions(options ...*Options) *DataBuilder {
	d.data.YOptions = append(d.data.YOptions, options...)
	return d
}

// Build .
func (d *DataBuilder) Build() *Data {
	return d.data
}

// New .
func New(title string) *Data {
	return &Data{
		Title:      title,
		Dimensions: *new([]string),
		XAxis:      new(Axis),
		YAxis:      *new([]*Axis),
		XOptions:   new(Options),
		YOptions:   *new([]*Options),
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

// SetXOptions .
func (d *Data) SetXOptions(options *Options) {
	d.XOptions = options
}

// SetYOptions .
func (d *Data) SetYOptions(options ...*Options) {
	d.YOptions = append(d.YOptions, options...)
}
