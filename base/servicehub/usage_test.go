// Author: recallsong
// Email: songruiguo@qq.com

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
			args: args{},
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
