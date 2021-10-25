package set

import (
	"testing"
)

func Test_syncSet_Add(t *testing.T) {

	type args struct {
		element interface{}
	}
	tests := []struct {
		name string
		set  Set
		args args
		want bool
	}{
		{"case add success", NewSyncSet(), args{"test_element"}, true},
		{"case add failed", NewSyncSet("test_element"), args{"test_element"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Add(tt.args.element); got != tt.want {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncSet_Remove(t *testing.T) {
	type args struct {
		element interface{}
	}
	tests := []struct {
		name string
		set  Set
		args args
		want bool
	}{
		{"case remove 1", NewSyncSet("e1", "e2"), args{"e2"}, false},
		{"case remove 2", NewSyncSet("e1", "e3"), args{"e2"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set.Remove(tt.args.element)
			if got := tt.set.Contains(tt.args.element); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncSet_Contains(t *testing.T) {
	type args struct {
		elements []interface{}
	}
	tests := []struct {
		name string
		set  Set
		args args
		want bool
	}{
		{"case contains all", NewSyncSet("e1", "e2"), args{elements: []interface{}{"e1", "e2"}}, true},
		{"case not contains", NewSyncSet("e1"), args{elements: []interface{}{"e1", "e2"}}, false},
		{"case contains one", NewSyncSet("e1", "e2"), args{elements: []interface{}{"e1"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Contains(tt.args.elements...); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_syncSet_Clear(t *testing.T) {

	tests := []struct {
		name       string
		set        Set
		wantLength int
		wantErr    bool
	}{
		{"case clear1", NewSyncSet("e1", "e2"), 0, false},
		{"case clear2", NewSyncSet("e1", "e2"), 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set.Clear()
			if got := tt.set.Len(); (got != tt.wantLength) != tt.wantErr {
				t.Errorf("Len() = %v, want %v", got, tt.wantLength)
			}
		})
	}
}

func Test_syncSet_Len(t *testing.T) {

	tests := []struct {
		name string
		set  Set
		want int
	}{
		{"case len 0", NewSyncSet(), 0},
		{"case len 1", NewSyncSet("e1"), 1},
		{"case len 2", NewSyncSet("e1", "e2"), 2},
		{"case len 3", NewSyncSet("e1", "e2", "e3"), 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
