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

package runtime

import (
	"reflect"
	"testing"
)

func TestCompile(t *testing.T) {
	type args struct {
		path string
		url  string
	}
	type want struct {
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
			args: args{path: "/", url: "/"},
			want: want{},
		},
		{
			name: "static path",
			args: args{path: "/abc/def/g", url: "/abc/def/g"},
			want: want{},
		},
		{
			name: "one path param",
			args: args{
				path: "/abc/{def}/g",
				url:  "/abc/123/g",
			},
			want: want{params: map[string]string{"def": "123"}},
		},
		{
			name: "two path params",
			args: args{
				path: "/abc/{def}/g/{xyz}",
				url:  "/abc/123/g/456",
			},
			want: want{params: map[string]string{"def": "123", "xyz": "456"}},
		},
		{
			name: "many path params",
			args: args{
				path: "/abc/{def}/g/{h}/{x}/{yz}",
				url:  "/abc/123/g/4/56/7",
			},
			want: want{params: map[string]string{"def": "123", "h": "4", "x": "56", "yz": "7"}},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx}/yz",
				url:  "/abc/123/g/xx/yz",
			},
			want: want{params: map[string]string{"def": "123", "h": "xx"}},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx}/yz",
				url:  "/abc/123/g/xx123/yz",
			},
			want: want{notmatch: true},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx/123/*}/yz",
				url:  "/abc/123/g/xx/123/456/yz",
			},
			want: want{params: map[string]string{"def": "123", "h": "xx/123/456"}},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx/123/*}/yz:aaa",
				url:  "/abc/123/g/xx/123/456/yz",
			},
			want: want{notmatch: true},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx/123/*}/yz:aaa",
				url:  "/abc/123/g/xx/123/456/yz:bbb",
			},
			want: want{notmatch: true},
		},
		{
			args: args{
				path: "/abc/{def}/g/{h=xx/123/*}/yz:aaa",
				url:  "/abc/123/g/xx/123/456/yz:aaa",
			},
			want: want{params: map[string]string{"def": "123", "h": "xx/123/456"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher, err := Compile(tt.args.path)
			if err != nil {
				t.Errorf("Compile() return %s", err)
				return
			}
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
		})
	}
}
