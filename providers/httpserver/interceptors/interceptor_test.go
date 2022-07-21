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

package interceptors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo"
)

func Test_judgeAnyEnable(t *testing.T) {
	getC := func(url string, header http.Header) echo.Context {
		e := echo.New()
		r, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			panic(err)
		}
		r.Header = header
		c := e.NewContext(r, nil)
		return c
	}

	type args struct {
		c                echo.Context
		enableFetchFuncs []EnableFetchFunc
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all and default is disabled",
			args: args{
				c:                getC("localhost", nil),
				enableFetchFuncs: nil,
			},
			want: false,
		},
		{
			name: "first one is enabled",
			args: args{
				c: getC("localhost", nil),
				enableFetchFuncs: []EnableFetchFunc{
					func(c echo.Context) bool { return true },
				},
			},
			want: true,
		},
		{
			name: "the last one is enabled",
			args: args{
				c: getC("localhost", nil),
				enableFetchFuncs: []EnableFetchFunc{
					func(c echo.Context) bool { return false },
					func(c echo.Context) bool { return false },
					func(c echo.Context) bool { return true },
				},
			},
			want: true,
		},
		{
			name: "all disabled, but defined at url query",
			args: args{
				c:                getC(fmt.Sprintf("localhost?%s", defaultEnableFlag), nil),
				enableFetchFuncs: nil,
			},
			want: true,
		},
		{
			name: "all disabled, but defined at header",
			args: args{
				c:                getC("localhost", http.Header{defaultEnableFlag: []string{}}),
				enableFetchFuncs: nil,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := judgeAnyEnable(tt.args.c, tt.args.enableFetchFuncs); got != tt.want {
				t.Errorf("judgeAnyEnable() = %v, want %v", got, tt.want)
			}
		})
	}
}
