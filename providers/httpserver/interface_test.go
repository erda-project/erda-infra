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

package httpserver

import (
	"fmt"
	"testing"
)

func TestWithPathFormat(t *testing.T) {
	tests := []struct {
		name   string
		format PathFormat
		want   *pathFormater
	}{
		{
			name:   "googleapis",
			format: PathFormatGoogleAPIs,
			want: &pathFormater{
				typ:    PathFormatGoogleAPIs,
				format: buildGoogleAPIsPath,
				parser: googleAPIsPathParamsInterceptor,
			},
		},
		{
			name:   "echo path",
			format: PathFormatEcho,
			want: &pathFormater{
				typ:    PathFormatEcho,
				format: buildEchoPath,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			if got, ok := WithPathFormat(tt.format).(*pathFormater); !ok || fmt.Sprint(*got) != fmt.Sprint(*tt.want) {
				t.Errorf("WithPathFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
