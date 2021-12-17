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

type OptionBase struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

type SelectOption OptionBase

type SelectOptionWithChildren struct {
	SelectOption
	Children []SelectOption `json:"children,omitempty"`
}

func NewSelectOption(label string, value string) *SelectOption {
	return &SelectOption{
		Label: label,
		Value: value,
	}
}

func NewSelectChildrenOption(label string, value string, children []SelectOption) *SelectOptionWithChildren {
	return &SelectOptionWithChildren{
		SelectOption: *NewSelectOption(label, value),
		Children:     children,
	}
}
