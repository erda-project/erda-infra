// Author: recallsong
// Email: songruiguo@qq.com

package example

import (
	"testing"

	"github.com/erda-project/erda-infra/base/servicehub"
)

func getService(t *testing.T) Interface {
	hub := servicehub.Run(&servicehub.RunOptions{
		Content: `
example-provider:
`})
	example, ok := hub.Service("example").(Interface)
	if !ok {
		t.Fatalf("example is not Interface")
	}
	return example
}

func Test_provider_Hello(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test1",
			args{
				"test",
			},
			"hello test",
		},
		{
			"test2",
			args{
				"song",
			},
			"hello song",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := getService(t)
			if got := s.Hello(tt.args.name); got != tt.want {
				t.Errorf("provider.Hello() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_provider_Add(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"test1",
			args{
				a: 1,
				b: 2,
			},
			3,
		},
		{
			"test2",
			args{
				a: 8,
				b: 2,
			},
			10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := getService(t)
			if got := s.Add(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("provider.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_provider_sub(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"test1",
			args{
				a: 10,
				b: 2,
			},
			8,
		},
		{
			"test2",
			args{
				a: 10,
				b: 4,
			},
			6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := getService(t)
			p := s.(*provider)
			if got := p.sub(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("provider.sub() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (p *provider) testOnlyFunc(a, b int) int {
	return a + b
}

func Test_provider_testOnlyFunc(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"test1",
			args{
				a: 5,
				b: 7,
			},
			12,
		},
		{
			"test2",
			args{
				a: 10,
				b: 14,
			},
			24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := getService(t)
			p := s.(*provider)
			if got := p.testOnlyFunc(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("provider.testOnlyFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
