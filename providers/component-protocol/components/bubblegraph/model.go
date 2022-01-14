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

package bubblegraph

import (
	structure "github.com/erda-project/erda-infra/providers/component-protocol/components/commodel/data-structure"
)

// Below is standard struct for bubble graph related.
type (
	// Data includes List.
	Data struct {
		Title    string     `json:"title"`
		List     []*Bubble  `json:"list"`
		XOptions *Options   `json:"xOptions"`
		YOptions []*Options `json:"yOptions"`
	}

	// Bubble .
	Bubble struct {
		X         *Axis       `json:"x"`
		Y         *Axis       `json:"y"`
		Size      *BubbleSize `json:"size"`
		Group     string      `json:"group"`
		Dimension string      `json:"dimension"`
	}

	// BubbleSize .
	BubbleSize struct {
		Value float64 `json:"value"`
	}

	// Axis .
	Axis struct {
		Value interface{} `json:"value"`
		Unit  string      `json:"unit"`
	}

	// Options .
	Options struct {
		Dimension string                   `json:"dimension,omitempty"`
		Structure *structure.DataStructure `json:"structure"`
	}

	// BubbleBuilder .
	BubbleBuilder struct {
		bubble *Bubble
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
	return &DataBuilder{data: &Data{XOptions: &Options{}, YOptions: []*Options{}}}
}

// WithTitle .
func (d *DataBuilder) WithTitle(title string) *DataBuilder {
	d.data.Title = title
	return d
}

// WithBubble .
func (d *DataBuilder) WithBubble(bubbles ...*Bubble) *DataBuilder {
	d.data.List = append(d.data.List, bubbles...)
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

// NewBubbleBuilder .
func NewBubbleBuilder() *BubbleBuilder {
	return &BubbleBuilder{bubble: &Bubble{X: &Axis{}, Y: &Axis{}, Size: &BubbleSize{}}}
}

// WithValueX .
func (bb *BubbleBuilder) WithValueX(v interface{}) *BubbleBuilder {
	bb.bubble.X.Value = v
	return bb
}

// WithX .
func (bb *BubbleBuilder) WithX(x *Axis) *BubbleBuilder {
	bb.bubble.X = x
	return bb
}

// WithValueY .
func (bb *BubbleBuilder) WithValueY(v interface{}) *BubbleBuilder {
	bb.bubble.Y.Value = v
	return bb
}

// WithY .
func (bb *BubbleBuilder) WithY(y *Axis) *BubbleBuilder {
	bb.bubble.Y = y
	return bb
}

// WithSize .
func (bb *BubbleBuilder) WithSize(size *BubbleSize) *BubbleBuilder {
	bb.bubble.Size = size
	return bb
}

// WithValueSize .
func (bb *BubbleBuilder) WithValueSize(v float64) *BubbleBuilder {
	bb.bubble.Size.Value = v
	return bb
}

// WithGroup .
func (bb *BubbleBuilder) WithGroup(group string) *BubbleBuilder {
	bb.bubble.Group = group
	return bb
}

// WithDimension .
func (bb *BubbleBuilder) WithDimension(dimension string) *BubbleBuilder {
	bb.bubble.Dimension = dimension
	return bb
}

// Build .
func (bb *BubbleBuilder) Build() *Bubble {
	return bb.bubble
}
