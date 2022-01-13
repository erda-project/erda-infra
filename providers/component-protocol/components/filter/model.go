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

package filter

import (
	model "github.com/erda-project/erda-infra/providers/component-protocol/components/filter/models"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// Data .
type (
	// Data filter std data
	Data struct {
		// models/condition.go define constructors of condition type
		// SelectCondition, DateRangeCondition, etc
		Conditions []interface{}                            `json:"conditions,omitempty"`
		FilterSet  []SetItem                                `json:"filterSet,omitempty"`
		Operations map[cptype.OperationKey]cptype.Operation `json:"operations,omitempty"`
		// HideSave hide saveFilterSet
		HideSave bool `json:"hideSave,omitempty"`
	}

	// SetItem custom filter conditions
	SetItem struct {
		ID       string          `json:"id,omitempty"`
		Values   cptype.ExtraMap `json:"values,omitempty"`
		Label    string          `json:"label,omitempty"`
		IsPreset bool            `json:"isPreset,omitempty"`
	}

	// ICondition get type ICondition
	ICondition interface {
		Type() model.ConditionType
	}

	// // State filter std state
	// State struct {
	// 	Values            cptype.ExtraMap `json:"values,omitempty"`
	// 	SelectedFilterSet string          `json:"selectedFilterSet,omitempty"`
	// 	cptype.ExtraMap
	// }
)
