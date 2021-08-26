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
	"reflect"
	"testing"

	"github.com/erda-project/erda-infra/pkg/transport/http/runtime"
	"github.com/erda-project/erda-infra/providers/httpserver/server"
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
		{
			args: args{
				path: "/abc/{_}/d",
				url:  "/abc/123/d",
			},
			want: want{
				path:   "/abc/:_/d",
				params: map[string]string{"_": "123"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildGoogleAPIsPath(tt.args.path); got != tt.want.path {
				t.Errorf("buildGoogleAPIsPath() = %v, want %v", got, tt.want.path)
				return
			}
			if tt.want.skip {
				return
			}
			matcher, err := runtime.Compile(tt.args.path)
			if err != nil {
				t.Errorf("runtime.Compile() return %s", err)
				return
			}
			if len(tt.args.url) > 0 {
				vars, err := matcher.Match(tt.args.url)
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

func Test_buildEchoPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{
			path: "/abc/def",
			want: "/abc/def",
		},
		{
			path: "/abc/:def/g",
			want: "/abc/:def/g",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := buildEchoPath(tt.path); got != tt.want {
				t.Errorf("buildEchoPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_googleAPIsPathParamsInterceptor(t *testing.T) {
	handler := func(c server.Context) error { return nil }
	type args struct {
		path    string
		handler server.HandlerFunc
	}
	tests := []struct {
		name   string
		args   args
		static bool
	}{
		{
			name: "static",
			args: args{
				path:    "/abc/def",
				handler: handler,
			},
			static: true,
		},
		{
			name: "params",
			args: args{
				path:    "/abc/{def}/g",
				handler: handler,
			},
			static: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := googleAPIsPathParamsInterceptor(tt.args.path)(tt.args.handler)
			if tt.static {
				if fmt.Sprint(handler) != fmt.Sprint(tt.args.handler) {
					t.Errorf("googleAPIsPathParamsInterceptor()(handler) get non static handler")
				}
			} else {
				if fmt.Sprint(handler) == fmt.Sprint(tt.args.handler) {
					t.Errorf("googleAPIsPathParamsInterceptor()(handler) got wrapped handler")
				}
			}
		})
	}
}

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
