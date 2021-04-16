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
	"reflect"
	"strings"
	"testing"

	"github.com/erda-project/erda-infra/pkg/transport/http/httprule"
	"github.com/erda-project/erda-infra/pkg/transport/http/runtime"
)

func Test_buildGoogleAPIsPath(t *testing.T) {
	type args struct {
		path string
		url  string
	}
	type want struct {
		path     string
		skip     bool
		params   map[string]string
		notmatch bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "empty path",
			args: args{},
			want: want{},
		},
		{
			name: "root path",
			args: args{path: "/"},
			want: want{path: "/"},
		},
		{
			name: "static path",
			args: args{path: "/abc/def/g"},
			want: want{path: "/abc/def/g"},
		},
		{
			name: "one path param",
			args: args{
				path: "/abc/{def}/g",
				url:  "/abc/123/g",
			},
			want: want{
				path:   "/abc/:def/g",
				params: map[string]string{"def": "123"},
			},
		},
		{
			name: "two path params",
			args: args{
				path: "/abc/{def}/g/{xyz}",
				url:  "/abc/123/g/456",
			},
			want: want{
				path:   "/abc/:def/g/:xyz",
				params: map[string]string{"def": "123", "xyz": "456"},
			},
		},
		{
			name: "many path params",
			args: args{
				path: "/abc/{def}/g/{h}/{x}/{yz}",
				url:  "/abc/123/g/4/56/7",
			},
			want: want{
				path:   "/abc/:def/g/:h/:x/:yz",
				params: map[string]string{"def": "123", "h": "4", "x": "56", "yz": "7"},
			},
		},
		{
			name: "has empty path param",
			args: args{path: "/abc/{def}/g/{}/{x}/{yz}"},
			want: want{path: "/abc/:def/g/*/:x/:yz", skip: true},
		},
		{
			name: "* path param",
			args: args{path: "/abc/{def}/g/{*}/{x}/{yz}"},
			want: want{path: "/abc/:def/g/*/:x/:yz", skip: true},
		},
		{
			args: args{path: "/abc/{def}/g/{h/{x}/{yz}"},
			want: want{path: "/abc/:def/g/:h%2F%7Bx/:yz", skip: true},
		},
		{
			args: args{path: "/abc/{def}/g/{h}/{x}/{yz"},
			want: want{path: "/abc/:def/g/:h/:x/{yz", skip: true},
		},
		{
			args: args{path: "/abc/{def}/g/{h=xx}/{yz"},
			want: want{path: "/abc/:def/g/:h/{yz", skip: true},
		},
		{
			args: args{path: "/abc/{def}/g/{h=xx}/{yz="},
			want: want{path: "/abc/:def/g/:h/{yz=", skip: true},
		},
		{
			args: args{path: "/abc/{def}/g/{h=xx/***}/{yz="},
			want: want{path: "/abc/:def/g/:h/***/{yz=", skip: true},
		},
		{
			args: args{path: "/abc/{def}/g/{h=xx/**}/yz:verb"},
			want: want{path: "/abc/:def/g/:h/**/yz%3Averb", skip: true},
		},
		{
			args: args{path: "/abc/{def}/g/{h=xx/**}/yz:verb1:verb2"},
			want: want{path: "/abc/:def/g/:h/**/yz%3Averb1%3Averb2", skip: true},
		},
		{
			name: "query string",
			args: args{path: "/abc/{def}/g/{h=xx/**}/yz:verb1:verb2?abc=123&def=456"},
			want: want{path: "/abc/:def/g/:h/**/yz%3Averb1%3Averb2", skip: true},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx}/yz",
				url:  "/abc/123/g/xx/yz",
			},
			want: want{
				path:   "/abc/:def/g/:h/yz",
				params: map[string]string{"def": "123", "h": "xx"},
			},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx}/yz",
				url:  "/abc/123/g/xx123/yz",
			},
			want: want{
				path:     "/abc/:def/g/:h/yz",
				notmatch: true,
			},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx/123/*}/yz",
				url:  "/abc/123/g/xx/123/456/yz",
			},
			want: want{
				path:   "/abc/:def/g/:h/123/*/yz",
				params: map[string]string{"def": "123", "h": "xx/123/456"},
			},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx/123/*}/yz:aaa",
				url:  "/abc/123/g/xx/123/456/yz",
			},
			want: want{
				path:     "/abc/:def/g/:h/123/*/yz%3Aaaa",
				notmatch: true,
			},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx/123/*}/yz:aaa",
				url:  "/abc/123/g/xx/123/456/yz:bbb",
			},
			want: want{
				path:     "/abc/:def/g/:h/123/*/yz%3Aaaa",
				notmatch: true,
			},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx/123/*}/yz:aaa",
				url:  "/abc/123/g/xx/123/456/yz:aaa",
			},
			want: want{
				path:   "/abc/:def/g/:h/123/*/yz%3Aaaa",
				params: map[string]string{"def": "123", "h": "xx/123/456"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildGoogleAPIsPath(tt.args.path); got != tt.want.path {
				t.Errorf("buildGoogleAPIsPath() = %v, want %v", got, tt.want.path)
				return
			}
			if len(tt.want.path) <= 0 || tt.want.path == "/" || tt.want.skip {
				return
			}
			compiler, err := httprule.Parse(tt.args.path)
			if err != nil {
				t.Errorf("httprule.Parse() return %s", err)
				return
			}
			temp := compiler.Compile()
			pattern, err := runtime.NewPattern(httprule.SupportPackageIsVersion1, temp.OpCodes, temp.Pool, temp.Verb)
			if err != nil {
				t.Errorf("runtime.NewPattern() return %s", err)
				return
			}
			if len(tt.args.url) > 0 {
				components := strings.Split(tt.args.url[1:], "/")
				last := len(components) - 1
				var verb string
				if idx := strings.LastIndex(components[last], ":"); idx >= 0 {
					c := components[last]
					components[last], verb = c[:idx], c[idx+1:]
				}
				vars, err := pattern.Match(components, verb)
				if err != nil {
					if !tt.want.notmatch {
						t.Errorf("not match %q by %q", tt.args.url, tt.args.path)
						return
					}
					return
				}
				if len(vars) == 0 && len(tt.want.params) == 0 {
					return
				}
				if !reflect.DeepEqual(vars, tt.want.params) {
					t.Errorf("params not match, got %v, want %v", vars, tt.want.params)
					return
				}
			}
		})
	}
}
