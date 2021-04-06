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

// Reference: https://github.com/grpc-ecosystem/grpc-gateway/blob/v2.3.0/internal/httprule/types_test.go

package httprule

import (
	"fmt"
	"testing"
)

func TestTemplateStringer(t *testing.T) {
	for _, spec := range []struct {
		segs []segment
		want string
	}{
		{
			segs: []segment{
				literal("v1"),
			},
			want: "/v1",
		},
		{
			segs: []segment{
				wildcard{},
			},
			want: "/*",
		},
		{
			segs: []segment{
				deepWildcard{},
			},
			want: "/**",
		},
		{
			segs: []segment{
				variable{
					path: "name",
					segments: []segment{
						literal("a"),
					},
				},
			},
			want: "/{name=a}",
		},
		{
			segs: []segment{
				variable{
					path: "name",
					segments: []segment{
						literal("a"),
						wildcard{},
						literal("b"),
					},
				},
			},
			want: "/{name=a/*/b}",
		},
		{
			segs: []segment{
				literal("v1"),
				variable{
					path: "name",
					segments: []segment{
						literal("a"),
						wildcard{},
						literal("b"),
					},
				},
				literal("c"),
				variable{
					path: "field.nested",
					segments: []segment{
						wildcard{},
						literal("d"),
					},
				},
				wildcard{},
				literal("e"),
				deepWildcard{},
			},
			want: "/v1/{name=a/*/b}/c/{field.nested=*/d}/*/e/**",
		},
	} {
		tmpl := template{segments: spec.segs}
		if got, want := tmpl.String(), spec.want; got != want {
			t.Errorf("%#v.String() = %q; want %q", tmpl, got, want)
		}

		tmpl.verb = "LOCK"
		if got, want := tmpl.String(), fmt.Sprintf("%s:LOCK", spec.want); got != want {
			t.Errorf("%#v.String() = %q; want %q", tmpl, got, want)
		}
	}
}
