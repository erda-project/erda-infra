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

package complexgraph

import structure "github.com/erda-project/erda-infra/providers/component-protocol/components/commodel/data-structure"

// Y axis positions
const (
	Right Position = "right"
	Left  Position = "left"
)

// X axis positions
const (
	Top    Position = "top"
	Bottom Position = "bottom"
)

// Define type
const (
	Value    Type = "value"
	Category Type = "category"
	Bar      Type = "bar"
	Line     Type = "line"
)

// Below is standard struct for graph related.
type (
	// Data includes List.
	Data struct {
		Title      string   `json:"title"`
		Dimensions []string `json:"dimensions"`
		XAxis      []*Axis  `json:"xAxis"`
		YAxis      []*Axis  `json:"yAxis"`
		Series     []*Sere  `json:"series"`
		Inverse    bool     `json:"inverse"`
	}

	//Axis .
	Axis struct {
		Type          Type                     `json:"type"`
		Name          string                   `json:"name,omitempty"`
		Dimensions    []string                 `json:"dimensions,omitempty"`
		Position      Position                 `json:"position,omitempty"`
		Data          []interface{}            `json:"data,omitempty"`
		DataStructure *structure.DataStructure `json:"structure"`
	}

	//Sere .
	Sere struct {
		Name      string        `json:"name"`
		Dimension string        `json:"dimension"`
		Type      Type          `json:"type"`
		Data      []interface{} `json:"data"`
	}

	//Type defined graph type
	Type string

	//Position defined axis position
	Position string

	// DataBuilder .
	DataBuilder struct {
		data *Data
	}

	// AxisBuilder .
	AxisBuilder struct {
		data *Axis
	}

	// SereBuilder .
	SereBuilder struct {
		data *Sere
	}
)

// NewSereBuilder .
func NewSereBuilder() *SereBuilder {
	return &SereBuilder{data: &Sere{}}
}

// Build .
func (d *SereBuilder) Build() *Sere {
	return d.data
}

// WithName .
func (d *SereBuilder) WithName(name string) *SereBuilder {
	d.data.Name = name
	return d
}

// WithDimension .
func (d *SereBuilder) WithDimension(dimension string) *SereBuilder {
	d.data.Dimension = dimension
	return d
}

// WithType .
func (d *SereBuilder) WithType(typ Type) *SereBuilder {
	d.data.Type = typ
	return d
}

// WithData .
func (d *SereBuilder) WithData(data ...interface{}) *SereBuilder {
	for _, datum := range data {
		d.data.Data = append(d.data.Data, datum)
	}
	return d
}

// NewAxisBuilder .
func NewAxisBuilder() *AxisBuilder {
	return &AxisBuilder{data: &Axis{DataStructure: &structure.DataStructure{}}}
}

// Build .
func (d *AxisBuilder) Build() *Axis {
	return d.data
}

// WithName .
func (d *AxisBuilder) WithName(name string) *AxisBuilder {
	d.data.Name = name
	return d
}

// WithDimensions .
func (d *AxisBuilder) WithDimensions(dimensions ...string) *AxisBuilder {
	for _, dimension := range dimensions {
		d.data.Dimensions = append(d.data.Dimensions, dimension)
	}
	return d
}

// WithPosition .
func (d *AxisBuilder) WithPosition(position Position) *AxisBuilder {
	d.data.Position = position
	return d
}

// WithData .
func (d *AxisBuilder) WithData(data ...interface{}) *AxisBuilder {
	for _, datum := range data {
		d.data.Data = append(d.data.Data, datum)
	}
	return d
}

// WithDataStructure .
func (d *AxisBuilder) WithDataStructure(typ structure.Type, precision structure.Precision, enable bool) *AxisBuilder {
	d.data.DataStructure.Type = typ
	d.data.DataStructure.Precision = precision
	d.data.DataStructure.Enable = enable
	return d
}

// WithType .
func (d *AxisBuilder) WithType(typ Type) *AxisBuilder {
	d.data.Type = typ
	return d
}

// NewDataBuilder .
func NewDataBuilder() *DataBuilder {
	return &DataBuilder{data: &Data{}}
}

// Build .
func (d *DataBuilder) Build() *Data {
	return d.data
}

// WithTitle .
func (d *DataBuilder) WithTitle(title string) *DataBuilder {
	d.data.Title = title
	return d
}

// WithDimensions .
func (d *DataBuilder) WithDimensions(dimensions ...string) *DataBuilder {
	for _, dimension := range dimensions {
		d.data.Dimensions = append(d.data.Dimensions, dimension)
	}
	return d
}

func (d *DataBuilder) EnableInverse() *DataBuilder {
	d.data.Inverse = true
	return d
}

// WithXAxis .
func (d *DataBuilder) WithXAxis(xAxis ...*Axis) *DataBuilder {
	for _, axis := range xAxis {
		d.data.XAxis = append(d.data.XAxis, axis)
	}
	return d
}

// WithYAxis .
func (d *DataBuilder) WithYAxis(yAxis ...*Axis) *DataBuilder {
	for _, axis := range yAxis {
		d.data.YAxis = append(d.data.YAxis, axis)
	}
	return d
}

// WithSeries .
func (d *DataBuilder) WithSeries(series ...*Sere) *DataBuilder {
	for _, sere := range series {
		d.data.Series = append(d.data.Series, sere)
	}
	return d
}
