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

package protocol

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

func Test_getCompNameAndInstanceName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name             string
		args             args
		wantCompName     string
		wantInstanceName string
	}{
		{
			name:             "with @",
			args:             args{name: "mt_block_detail_item@mt_case_num_total"},
			wantCompName:     "mt_block_detail_item",
			wantInstanceName: "mt_case_num_total",
		},
		{
			name:             "without @",
			args:             args{name: "mt_case_num_total"},
			wantCompName:     "mt_case_num_total",
			wantInstanceName: "mt_case_num_total",
		},
		{
			name:             "with @@",
			args:             args{name: "a@@b"},
			wantCompName:     "a",
			wantInstanceName: "@b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCompName, gotInstanceName := getCompNameAndInstanceName(tt.args.name)
			if gotCompName != tt.wantCompName {
				t.Errorf("getCompNameAndInstanceName() gotCompName = %v, want %v", gotCompName, tt.wantCompName)
			}
			if gotInstanceName != tt.wantInstanceName {
				t.Errorf("getCompNameAndInstanceName() gotInstanceName = %v, want %v", gotInstanceName, tt.wantInstanceName)
			}
		})
	}
}

func Test_calculateDefaultRenderOrderByHierarchy(t *testing.T) {
	py := `
hierarchy:
  root: page
  structure:
    page:
      - filter
      - overview_group
      - blocks
      - mt_plan_chart_group
      - leftHead
    overview_group:
      - quality_chart
      - blocks
    blocks:
      - mt_block
      - at_block
    mt_block:
      - mt_block_header
      - mt_block_detail
    mt_block_header:
      right: mt_block_header_filter
      left: mt_block_header_title
    mt_plan_chart_group:
      children:
        - mt_plan_chart2
        - mt_plan_chart
      extraContent: mt_plan_chart_filter
    leftHead:
      right:
        - leftHeadAddSceneSet
        - moreOperation
      left: leftHeadTitle
      tabBarExtraContent:
        - tabSceneSetExecuteButton
`
	var p cptype.ComponentProtocol
	assert.NoError(t, yaml.Unmarshal([]byte(py), &p))
	orders, err := calculateDefaultRenderOrderByHierarchy(&p)
	assert.NoError(t, err)
	expected := []string{"page", "filter", "overview_group", "quality_chart", "blocks", "mt_block",
		"mt_block_header", "mt_block_header_title", "mt_block_header_filter",
		"mt_block_detail", "at_block",
		"mt_plan_chart_group", "mt_plan_chart_filter", "mt_plan_chart2", "mt_plan_chart",
		"leftHead", "leftHeadTitle", "leftHeadAddSceneSet", "moreOperation", "tabSceneSetExecuteButton",
	}
	assert.Equal(t, expected, orders)
}

func Test_recursiveWalkCompOrder(t *testing.T) {
	// recursive walk from root
	var result []string
	allCompSubMap := make(map[string][]string)
	allCompSubMap["page"] = []string{"title", "overview", "filter"}
	allCompSubMap["overview"] = []string{"quality_chart", "blocks"}
	err := recursiveWalkCompOrder("page", &result, allCompSubMap)
	assert.NoError(t, err)
	expected := []string{"page", "title", "overview", "quality_chart", "blocks", "filter"}
	assert.Equal(t, expected, result)
}

func Test_getDefaultHierarchyRenderOrderFromCompExclude(t *testing.T) {
	type args struct {
		fullOrders           []string
		startFromCompExclude string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "not found, return nil",
			args: args{
				fullOrders:           []string{"root", "page", "overview", "split"},
				startFromCompExclude: "not",
			},
			want: []string{},
		},
		{
			name: "found in center",
			args: args{
				fullOrders:           []string{"root", "page", "overview", "split"},
				startFromCompExclude: "page",
			},
			want: []string{"overview", "split"},
		},
		{
			name: "found in the beginning",
			args: args{
				fullOrders:           []string{"root", "page", "overview", "split"},
				startFromCompExclude: "root",
			},
			want: []string{"page", "overview", "split"},
		},
		{
			name: "found in the last",
			args: args{
				fullOrders:           []string{"root", "page", "overview", "split"},
				startFromCompExclude: "split",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDefaultHierarchyRenderOrderFromCompExclude(tt.args.fullOrders, tt.args.startFromCompExclude); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefaultHierarchyRenderOrderFromCompExclude() = %v, want %v", got, tt.want)
			}
		})
	}
}
