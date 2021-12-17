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

type ConditionType string

const (
	ConditionTypeSelect      ConditionType = "select"
	ConditionTypeInput       ConditionType = "input"
	ConditionTypeDateRange   ConditionType = "dateRange"
	ConditionTypeRangePicker ConditionType = "rangePicker"
)

type ConditionBase struct {
	Key         string        `json:"key,omitempty"`
	Label       string        `json:"label,omitempty"`
	Type        ConditionType `json:"type,omitempty"`
	Placeholder string        `json:"placeholder,omitempty"`
}

type SelectCondition struct {
	ConditionBase
	Mode    string         `json:"mode,omitempty"`
	Options []SelectOption `json:"options,omitempty"`
}

type SelectConditionWithChildren struct {
	ConditionBase
	Mode    string                     `json:"mode,omitempty"`
	Options []SelectOptionWithChildren `json:"options,omitempty"`
}

type DateRangeCondition SelectCondition

func (o SelectCondition) Type() ConditionType {
	return ConditionTypeSelect
}

func (o SelectConditionWithChildren) Type() ConditionType {
	return ConditionTypeSelect
}

func (o DateRangeCondition) Type() ConditionType {
	return ConditionTypeDateRange
}

func NewCondition(key string, label string) *ConditionBase {
	return &ConditionBase{
		Key:   key,
		Label: label,
	}
}

// initial condition with select option
func NewSelectCondition(key string, label string, options []SelectOption) *SelectCondition {
	var r = SelectCondition{
		ConditionBase: *NewCondition(key, label),
		Options:       options,
	}
	r.ConditionBase.Type = r.Type()
	return &r
}

// initial condition with select option with children
func NewSelectConditionWithChildren(key string, label string, options []SelectOptionWithChildren) *SelectConditionWithChildren {
	var r = SelectConditionWithChildren{
		ConditionBase: *NewCondition(key, label),
		Options:       options,
	}
	r.ConditionBase.Type = r.Type()
	return &r
}

func NewDateRangeCondition(key string, label string) *DateRangeCondition {
	var r = DateRangeCondition{
		ConditionBase: *NewCondition(key, label),
	}
	r.ConditionBase.Type = r.Type()
	return &r
}

func (s *SelectCondition) WithPlaceHolder(placeholder string) *SelectCondition {
	s.Placeholder = placeholder
	return s
}

func (s *SelectCondition) WithMode(mode string) *SelectCondition {
	s.Mode = mode
	return s
}