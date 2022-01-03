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

package cputil

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go/cp/pb"
)

func TestObjJSONTransfer(t *testing.T) {
	src := pb.ComponentProtocol{
		Options: &pb.ProtocolOptions{
			SyncIntervalSecond: 0.01,
		},
		Hierarchy: &pb.Hierarchy{
			Parallel: map[string]*structpb.Value{
				"page": func() *structpb.Value {
					result, err := structpb.NewValue([]interface{}{"filter", "grid"})
					if err != nil {
						panic(err)
					}
					return result
				}(),
			},
		},
	}
	var dest cptype.ComponentProtocol

	err := ObjJSONTransfer(&src, &dest)
	assert.NoError(t, err)
	assert.Equal(t, src.Options.SyncIntervalSecond, dest.Options.SyncIntervalSecond)
	fmt.Printf("%#v\n", dest.Hierarchy.Parallel)
}

func TestMustFlatMapMeta(t *testing.T) {
	input := map[string]interface{}{
		"a": "b",
	}
	MustFlatMapMeta(input, false)
	assert.True(t, len(input) == 1)

	input = map[string]interface{}{
		"c": "d",
		"meta": map[string]interface{}{
			"a": "b",
		},
	}
	MustFlatMapMeta(input, false)
	assert.True(t, len(input) == 3)

	// map
	input = map[string]interface{}{
		"c": "d",
		"meta": map[string]interface{}{
			"a": "b",
		},
		"flatMeta": true,
	}
	MustFlatMapMeta(input, false)
	assert.True(t, len(input) == 4)
	assert.Equal(t, "b", input["a"])

	// ref map
	input = map[string]interface{}{
		"c": "d",
		"meta": map[string]interface{}{
			"a": "b",
			"e": "f",
		},
		"flatMeta": true,
	}
	MustFlatMapMeta(&input, false)
	assert.True(t, len(input) == 5)
	assert.Equal(t, "b", input["a"])
	assert.Equal(t, "f", input["e"])

	// nested flat meta
	input = map[string]interface{}{
		"a": "b",
		"meta": map[string]interface{}{
			"sub": map[string]interface{}{
				"meta": map[string]interface{}{
					"e": "f",
					"g": "h",
				},
			},
		},
	}
	MustFlatMapMeta(&input, true)
	b, _ := json.MarshalIndent(input, "", "  ")
	fmt.Println(string(b))
}
