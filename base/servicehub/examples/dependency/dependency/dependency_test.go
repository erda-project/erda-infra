// Author: recallsong
// Email: songruiguo@qq.com

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
