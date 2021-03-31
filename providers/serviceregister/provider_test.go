// Author: recallsong
// Email: songruiguo@qq.com

package register

import "testing"

func Test_buildPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "empty path",
			path: "",
			want: "",
		},
		{
			name: "root path",
			path: "/",
			want: "/",
		},
		{
			name: "static path",
			path: "/abc/def/g",
			want: "/abc/def/g",
		},
		{
			name: "one path param",
			path: "/abc/{def}/g",
			want: "/abc/:def/g",
		},
		{
			name: "two path params",
			path: "/abc/{def}/g/{xyz}",
			want: "/abc/:def/g/:xyz",
		},
		{
			name: "many path params",
			path: "/abc/{def}/g/{h}/{x}/{yz}",
			want: "/abc/:def/g/:h/:x/:yz",
		},
		{
			name: "has empty path param",
			path: "/abc/{def}/g/{}/{x}/{yz}",
			want: "/abc/:def/g/*/:x/:yz",
		},
		{
			name: "* path param",
			path: "/abc/{def}/g/{*}/{x}/{yz}",
			want: "/abc/:def/g/*/:x/:yz",
		},
		{
			path: "/abc/{def}/g/{h/{x}/{yz}",
			want: "/abc/:def/g/:h%2F%7Bx/:yz",
		},
		{
			path: "/abc/{def}/g/{h}/{x}/{yz",
			want: "/abc/:def/g/:h/:x/{yz",
		},
		{
			path: "/abc/{def}/g/{h=/xx}/{yz",
			want: "/abc/:def/g/:h/xx/{yz",
		},
		{
			path: "/abc/{def}/g/{h=/xx}/{yz=",
			want: "/abc/:def/g/:h/xx/{yz=",
		},
		{
			path: "/abc/{def}/g/{h=/xx/***}/{yz=",
			want: "/abc/:def/g/:h/xx/***/{yz=",
		},
		{
			path: "/abc/{def}/g/{h=/xx/**}/yz:verb",
			want: "/abc/:def/g/:h/xx/**/yz%3Averb",
		},
		{
			path: "/abc/{def}/g/{h=/xx/**}/yz:verb1:verb2",
			want: "/abc/:def/g/:h/xx/**/yz%3Averb1%3Averb2",
		},
		{
			name: "query string",
			path: "/abc/{def}/g/{h=/xx/**}/yz:verb1:verb2?abc=123&def=456",
			want: "/abc/:def/g/:h/xx/**/yz%3Averb1%3Averb2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildPath(tt.path); got != tt.want {
				t.Errorf("buildPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
