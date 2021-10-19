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
			args:             args{ name: "mt_block_detail_item@mt_case_num_total" },
			wantCompName:     "mt_block_detail_item",
			wantInstanceName: "mt_case_num_total",
		},
		{
			name:             "without @",
			args:             args{ name: "mt_case_num_total"},
			wantCompName:     "mt_case_num_total",
			wantInstanceName: "mt_case_num_total",
		},
		{
			name:             "with @@",
			args:             args{ name: "a@@b"},
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
`
	var p cptype.ComponentProtocol
	assert.NoError(t, yaml.Unmarshal([]byte(py), &p))
	orders := calculateDefaultRenderOrderByHierarchy(&p)
	expected := []string{"page", "filter", "overview_group", "quality_chart", "blocks", "mt_block",
		"mt_block_header", "mt_block_header_title", "mt_block_header_filter",
		"mt_block_detail", "at_block"}
	for i := range expected {
		assert.True(t,true, expected[i] == orders[i])
	}
}

func Test_recursiveWalkCompOrder(t *testing.T) {
	// recursive walk from root
	var result []string
	allCompSubMap := make(map[string][]string)
	allCompSubMap["page"] = []string{"title", "overview", "filter"}
	allCompSubMap["overview"] = []string{"quality_chart", "blocks"}
	recursiveWalkCompOrder("page", &result, allCompSubMap)
	expected := []string{"page", "title", "overview", "quality_chart", "blocks", "filter"}
	for i := range expected {
		assert.True(t,true, expected[i] == result[i])
	}
}