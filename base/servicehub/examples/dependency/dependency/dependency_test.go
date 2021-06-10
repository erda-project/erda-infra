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

package dependency

import (
	"testing"

	"github.com/erda-project/erda-infra/base/servicehub"
)

func Test_provider_Hello(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args{
				"test",
			},
			"hello test",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			hub := servicehub.Run(&servicehub.RunOptions{
				Content: `
example-dependency-provider:
`})
			hello, ok := hub.Service("example-dependency").(Interface)
			if !ok {
				t.Fatalf("example-dependency is not Interface")
			}
			if got := hello.Hello(tt.args.name); got != tt.want {
				t.Errorf("provider.Hello() = %v, want %v", got, tt.want)
			}
		})
	}
}
