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

package model

// ConditionType .
type ConditionType string

// ConditionTypeSelect .
const (
	ConditionTypeSelect     ConditionType = "select"
	ConditionTypeInput      ConditionType = "input"
	ConditionTypeDateRange  ConditionType = "dateRange"
	ConditionTypeTagsSelect ConditionType = "tagsSelect"
)

// ConditionBase .
type ConditionBase struct {
	Key         string        `json:"key,omitempty"`
	Label       string        `json:"label,omitempty"`
	Type        ConditionType `json:"type,omitempty"`
	Placeholder string        `json:"placeholder,omitempty"`
	Disabled    bool          `json:"disabled,omitempty"`
}

// SelectCondition .
type SelectCondition struct {
	ConditionBase
	Mode    string         `json:"mode,omitempty"`
	Options []SelectOption `json:"options,omitempty"`
}

// SelectConditionWithChildren .
type SelectConditionWithChildren struct {
	ConditionBase
	Mode    string                     `json:"mode,omitempty"`
	Options []SelectOptionWithChildren `json:"options,omitempty"`
}

// DateRangeCondition .
type DateRangeCondition SelectCondition

// TagsSelectCondition .
type TagsSelectCondition struct {
	ConditionBase
	Mode    string             `json:"mode,omitempty"`
	Options []TagsSelectOption `json:"options,omitempty"`
}

// Type .
func (o SelectCondition) Type() ConditionType {
	return ConditionTypeSelect
}

// Type .
func (o SelectConditionWithChildren) Type() ConditionType {
	return ConditionTypeSelect
}

// Type .
func (o DateRangeCondition) Type() ConditionType {
	return ConditionTypeDateRange
}

// Type .
func (o TagsSelectCondition) Type() ConditionType {
	return ConditionTypeTagsSelect
}

// NewCondition .
func NewCondition(key string, label string) *ConditionBase {
	return &ConditionBase{
		Key:   key,
		Label: label,
	}
}

// NewSelectCondition initial condition with select option
func NewSelectCondition(key string, label string, options []SelectOption) *SelectCondition {
	var r = SelectCondition{
		ConditionBase: *NewCondition(key, label),
		Options:       options,
	}
	r.ConditionBase.Type = r.Type()
	return &r
}

// NewSelectConditionWithChildren initial condition with select option with children
func NewSelectConditionWithChildren(key string, label string, options []SelectOptionWithChildren) *SelectConditionWithChildren {
	var r = SelectConditionWithChildren{
		ConditionBase: *NewCondition(key, label),
		Options:       options,
	}
	r.ConditionBase.Type = r.Type()
	return &r
}

// NewDateRangeCondition .
func NewDateRangeCondition(key string, label string) *DateRangeCondition {
	var r = DateRangeCondition{
		ConditionBase: *NewCondition(key, label),
	}
	r.ConditionBase.Type = r.Type()
	return &r
}

// NewTagsSelectCondition .
func NewTagsSelectCondition(key string, label string, options []TagsSelectOption) *TagsSelectCondition {
	var r = TagsSelectCondition{
		ConditionBase: *NewCondition(key, label),
		Options:       options,
	}
	r.ConditionBase.Type = r.Type()
	return &r
}

// WithPlaceHolder .
func (o *SelectCondition) WithPlaceHolder(placeholder string) *SelectCondition {
	o.Placeholder = placeholder
	return o
}

// WithMode .
func (o *SelectCondition) WithMode(mode string) *SelectCondition {
	o.Mode = mode
	return o
}
