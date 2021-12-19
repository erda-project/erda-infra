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

// OptionBase .
type OptionBase struct {
	Label string      `json:"label,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// SelectOption .
type SelectOption OptionBase

// SelectOptionWithChildren .
type SelectOptionWithChildren struct {
	SelectOption
	Children []SelectOption `json:"children,omitempty"`
}

// TagsSelectOption .
type TagsSelectOption struct {
	SelectOption
	Color string `json:"color,omitempty"`
}

// NewSelectOption .
func NewSelectOption(label string, value interface{}) *SelectOption {
	return &SelectOption{
		Label: label,
		Value: value,
	}
}

// NewSelectChildrenOption .
func NewSelectChildrenOption(label string, value interface{}, children []SelectOption) *SelectOptionWithChildren {
	return &SelectOptionWithChildren{
		SelectOption: *NewSelectOption(label, value),
		Children:     children,
	}
}

// NewTagsSelectOption .
func NewTagsSelectOption(label string, value interface{}, color string) *TagsSelectOption {
	return &TagsSelectOption{
		SelectOption: *NewSelectOption(label, value),
		Color:        color,
	}
}
