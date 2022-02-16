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
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

func TestGetUrlQuery(t *testing.T) {
	type FilterUrlQueryStruct struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:"age,omitempty"`
	}

	type args struct {
		sdk             *cptype.SDK
		resultStructPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil inParams",
			args: args{
				sdk:             &cptype.SDK{InParams: nil},
				resultStructPtr: nil,
			},
			wantErr: false,
		},
		{
			name: "inParams not nil, but url query is empty",
			args: args{
				sdk: &cptype.SDK{
					InParams: cptype.InParams{MakeCompUrlQueryKey("c1"): ""},
					Comp:     &cptype.Component{Name: "c1"},
				},
				resultStructPtr: nil,
			},
			wantErr: true,
		},
		{
			name: "normal",
			args: args{
				sdk: &cptype.SDK{
					InParams: cptype.InParams{
						MakeCompUrlQueryKey("c1"): func() string {
							q := FilterUrlQueryStruct{Name: "bob", Age: 20}
							b, err := json.Marshal(&q)
							assert.NoError(t, err)
							base64UrlQuery := base64.URLEncoding.EncodeToString(b)
							return base64UrlQuery
						}(),
					},
					Comp: &cptype.Component{Name: "c1"},
				},
				resultStructPtr: &FilterUrlQueryStruct{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetUrlQuery(tt.args.sdk, tt.args.resultStructPtr)
			if err != nil {
				t.Log(err)
			}
			assert.Equal(t, err != nil, tt.wantErr)
			fmt.Println(tt.args.resultStructPtr)
		})
	}
}

func TestEmptyBase64Decode(t *testing.T) {
	_, err := base64.StdEncoding.DecodeString("")
	assert.NoError(t, err)
}
