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

package servicehub

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func Test_providerContext_BindConfig(t *testing.T) {
	type config struct {
		StrVal             string
		IntVal             int
		Int8Val            int8
		Int16Val           int16
		Int32Val           int32
		Int64Val           int64
		BoolVal            bool
		DefaultVal         string `default:"default-val"`
		DurationVal        time.Duration
		DefaultDurationVal time.Duration `default:"3s"`
		Name               string        `file:"rename_field"`
		Name1Name2Name3    string
		DefaultVal2        string
	}
	tests := []struct {
		name    string
		pc      *providerContext
		flags   *pflag.FlagSet
		wantErr bool
		want    interface{}
	}{
		{
			name: "test all",
			pc: &providerContext{
				cfg: map[string]interface{}{
					"strval":          "test-val",
					"intval":          64,
					"int8val":         8,
					"int16val":        16,
					"int32val":        32,
					"int64val":        64,
					"boolval":         true,
					"durationval":     "5s",
					"rename_field":    "rename-val",
					"name1Name2Name3": "long-name-val",
				},
				define: &specDefine{&Spec{
					ConfigFunc: func() interface{} {
						return &config{
							DefaultVal2: "default-val-2",
						}
					},
				}},
			},
			want: &config{
				StrVal:             "test-val",
				IntVal:             64,
				Int8Val:            8,
				Int16Val:           16,
				Int32Val:           32,
				Int64Val:           64,
				BoolVal:            true,
				DefaultVal:         "default-val",
				DurationVal:        5 * time.Second,
				DefaultDurationVal: 3 * time.Second,
				Name:               "rename-val",
				Name1Name2Name3:    "long-name-val",
				DefaultVal2:        "default-val-2",
			},
		},
		{
			name: "test nil",
			pc: &providerContext{
				cfg: map[string]interface{}{
					"some-field": "test-val",
				},
				define: &specDefine{&Spec{
					ConfigFunc: func() interface{} {
						return nil
					},
				}},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pc.BindConfig(tt.flags); (err != nil) != tt.wantErr {
				t.Errorf("providerContext.BindConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.pc.cfg == nil && tt.want == nil {
				return
			}
			if !reflect.DeepEqual(tt.pc.cfg, tt.want) {
				t.Errorf("providerContext.cfg = %v, want %v", tt.pc.cfg, tt.want)
			}
		})
	}
}
