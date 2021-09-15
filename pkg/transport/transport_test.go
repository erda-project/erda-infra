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

package transport

import (
	"context"
	"reflect"
	"testing"

	transgrpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	transhttp "github.com/erda-project/erda-infra/pkg/transport/http"
	"github.com/erda-project/erda-infra/pkg/transport/interceptor"
)

type testInterceptor struct {
	key    string
	append func(v string)
}

func (i testInterceptor) Wrap(h interceptor.Handler) interceptor.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		i.append(i.key + "->")
		h(ctx, req)
		i.append("<-" + i.key)
		return nil, nil
	}
}

func TestWithInterceptors(t *testing.T) {
	tests := []struct {
		name   string
		handle string
		inters [][]*testInterceptor
		want   []string
	}{
		{
			handle: "handle",
			inters: [][]*testInterceptor{
				{{key: "a"}},
				{{key: "b"}},
				{{key: "c"}},
			},
			want: []string{"a->", "b->", "c->", "handle", "<-c", "<-b", "<-a"},
		},
		{
			handle: "handle",
			inters: [][]*testInterceptor{
				{{key: "a"}, {key: "b"}},
				{{key: "c"}},
			},
			want: []string{"a->", "b->", "c->", "handle", "<-c", "<-b", "<-a"},
		},
		{
			handle: "handle",
			inters: [][]*testInterceptor{
				{{key: "a"}, {key: "b"}},
				{{key: "c"}},
				{{key: "d"}, {key: "e"}},
			},
			want: []string{"a->", "b->", "c->", "d->", "e->", "handle", "<-e", "<-d", "<-c", "<-b", "<-a"},
		},
		{
			handle: "handle",
			want:   []string{"handle"},
		},
		{
			handle: "handle",
			inters: [][]*testInterceptor{
				{{key: "a"}},
			},
			want: []string{"a->", "handle", "<-a"},
		},
		{
			handle: "handle",
			inters: [][]*testInterceptor{
				{{key: "a"}},
				{{key: "b"}},
			},
			want: []string{"a->", "b->", "handle", "<-b", "<-a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getOpts := func(add func(v string)) (*ServiceOptions, interceptor.Handler) {
				handler := func(ctx context.Context, req interface{}) (interface{}, error) {
					add(tt.handle)
					return nil, nil
				}
				opts := DefaultServiceOptions()
				for _, inters := range tt.inters {
					var list []interceptor.Interceptor
					for _, i := range inters {
						i.append = add
						list = append(list, i.Wrap)
					}
					WithInterceptors(list...)(opts)
				}
				return opts, handler
			}
			var grpcResults []string
			opts, handler := getOpts(func(v string) {
				grpcResults = append(grpcResults, v)
			})
			grpcOpts := transgrpc.DefaultHandleOptions()
			for _, opt := range opts.GRPC {
				opt(grpcOpts)
			}
			if grpcOpts.Interceptor != nil {
				handler = grpcOpts.Interceptor(handler)
			}
			handler(nil, nil)
			if !reflect.DeepEqual(grpcResults, tt.want) {
				t.Errorf("wrapped grpc handler got %v, want %v", grpcResults, tt.want)
			}

			var httpResults []string
			opts, handler = getOpts(func(v string) {
				httpResults = append(httpResults, v)
			})
			httpOpts := transhttp.DefaultHandleOptions()
			for _, opt := range opts.HTTP {
				opt(httpOpts)
			}
			if httpOpts.Interceptor != nil {
				handler = httpOpts.Interceptor(handler)
			}
			handler(nil, nil)
			if !reflect.DeepEqual(httpResults, tt.want) {
				t.Errorf("wrapped http handler got %v, want %v", httpResults, tt.want)
			}
		})
	}
}
