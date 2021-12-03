package logrusx

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNew(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
	}{
		{"case", args{options: []Option{&option{logrus.InfoLevel}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.options...)
			got.Info("This will go to stdout")
			got.Warn("This will go to stderr")
		})
	}
}
