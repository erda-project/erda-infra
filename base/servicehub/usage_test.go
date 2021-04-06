// Copyright 2021 Terminus
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



package servicehub

import "testing"

func TestUsage(t *testing.T) {

	Register("test1-provider", &Spec{
		Services:    []string{"test"},
		Description: "this is provider for test1",
		ConfigFunc: func() interface{} {
			return &struct {
				Message string `file:"message" flag:"msg" default:"hi" desc:"message to show" env:"TEST_MESSAGE"`
			}{}
		},
		Creator: func() Provider {
			return &struct{}{}
		},
	})

	Register("test2-provider", &Spec{
		Services:    []string{"test"},
		Description: "this is provider for test2",
		ConfigFunc: func() interface{} {
			return &struct {
				Name string `file:"name" flag:"name" default:"test" desc:"description for test" env:"TEST_NAME"`
			}{}
		},
		Creator: func() Provider {
			return &struct{}{}
		},
	})

	type args struct {
		names []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				names: []string{"test1-provider"},
			},
			want: `Service Providers:
test1-provider
    this is provider for test1
    file:"message" flag:"msg" env:"TEST_MESSAGE" default:"hi" , message to show 
`,
		},
		{
			name: "test2",
			args: args{
				names: []string{"test2-provider"},
			},
			want: `Service Providers:
test2-provider
    this is provider for test2
    file:"name" flag:"name" env:"TEST_NAME" default:"test" , description for test 
`,
		},
		{
			name: "all providers",
			args: args{
				names: []string{"test1-provider", "test2-provider"},
			},
			want: `Service Providers:
test1-provider
    this is provider for test1
    file:"message" flag:"msg" env:"TEST_MESSAGE" default:"hi" , message to show 
test2-provider
    this is provider for test2
    file:"name" flag:"name" env:"TEST_NAME" default:"test" , description for test 
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Usage(tt.args.names...); got != tt.want {
				t.Errorf("Usage() = %v, want %v", got, tt.want)
			}
		})
	}
}
