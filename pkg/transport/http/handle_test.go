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

package http

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/erda-project/erda-infra/pkg/transport/interceptor"
)

type testInterceptor struct {
	key    string
	append func(v string)
}

func (i testInterceptor) WrapHTTP(h http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		i.append(i.key + "->")
		h(rw, r)
		i.append("<-" + i.key)
	}
}

func (i testInterceptor) Wrap(h interceptor.Handler) interceptor.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		i.append(i.key + "->")
		h(ctx, req)
		i.append("<-" + i.key)
		return nil, nil
	}
}

func TestWithHTTPInterceptor(t *testing.T) {
	tests := []struct {
		name   string
		handle string
		inters []*testInterceptor
		want   []string
	}{
		{
			handle: "handle",
			inters: []*testInterceptor{
				{key: "a"},
				{key: "b"},
				{key: "c"},
			},
			want: []string{"a->", "b->", "c->", "handle", "<-c", "<-b", "<-a"},
		},
		{
			handle: "handle",
			want:   []string{"handle"},
		},
		{
			handle: "handle",
			inters: []*testInterceptor{
				{key: "a"},
			},
			want: []string{"a->", "handle", "<-a"},
		},
		{
			handle: "handle",
			inters: []*testInterceptor{
				{key: "a"},
				{key: "b"},
			},
			want: []string{"a->", "b->", "handle", "<-b", "<-a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []string
			add := func(v string) {
				results = append(results, v)
			}
			handler := func(rw http.ResponseWriter, r *http.Request) {
				add(tt.handle)
			}
			opts := DefaultHandleOptions()
			for _, i := range tt.inters {
				i.append = add
				WithHTTPInterceptor(i.WrapHTTP)(opts)
			}
			if opts.HTTPInterceptor != nil {
				handler = opts.HTTPInterceptor(handler)
			}
			handler(nil, nil)
			if !reflect.DeepEqual(results, tt.want) {
				t.Errorf("wrapped http handler got %v, want %v", results, tt.want)
			}
		})
	}
}

func TestWithInterceptor(t *testing.T) {
	tests := []struct {
		name   string
		handle string
		inters []*testInterceptor
		want   []string
	}{
		{
			handle: "handle",
			inters: []*testInterceptor{
				{key: "a"},
				{key: "b"},
				{key: "c"},
			},
			want: []string{"a->", "b->", "c->", "handle", "<-c", "<-b", "<-a"},
		},
		{
			handle: "handle",
			want:   []string{"handle"},
		},
		{
			handle: "handle",
			inters: []*testInterceptor{
				{key: "a"},
			},
			want: []string{"a->", "handle", "<-a"},
		},
		{
			handle: "handle",
			inters: []*testInterceptor{
				{key: "a"},
				{key: "b"},
			},
			want: []string{"a->", "b->", "handle", "<-b", "<-a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var results []string
			add := func(v string) {
				results = append(results, v)
			}
			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				add(tt.handle)
				return nil, nil
			}
			opts := DefaultHandleOptions()
			for _, i := range tt.inters {
				i.append = add
				WithInterceptor(i.Wrap)(opts)
			}
			if opts.Interceptor != nil {
				handler = opts.Interceptor(handler)
			}
			handler(nil, nil)
			if !reflect.DeepEqual(results, tt.want) {
				t.Errorf("wrapped handler got %v, want %v", results, tt.want)
			}
		})
	}
}
