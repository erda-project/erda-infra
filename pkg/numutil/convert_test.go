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

package numutil

import (
	"math"
	"testing"
)

func TestMustInt(t *testing.T) {
	tests := []struct {
		name      string
		input     uint64
		expected  int
		wantPanic bool
	}{
		{"valid small value", 100, 100, false},
		{"valid zero", 0, 0, false},
		{"valid max int", uint64(math.MaxInt), math.MaxInt, false},
		{"overflow", uint64(math.MaxInt) + 1, 0, true},
		{"max uint64", math.MaxUint64, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("MustInt() expected to panic but didn't")
					}
				}()
			}
			result := MustInt(tt.input)
			if !tt.wantPanic && result != tt.expected {
				t.Errorf("MustInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMustInt64(t *testing.T) {
	tests := []struct {
		name      string
		input     uint64
		expected  int64
		wantPanic bool
	}{
		{"valid small value", 100, 100, false},
		{"valid zero", 0, 0, false},
		{"valid max int64", uint64(math.MaxInt64), math.MaxInt64, false},
		{"overflow", uint64(math.MaxInt64) + 1, 0, true},
		{"max uint64", math.MaxUint64, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("MustInt64() expected to panic but didn't")
					}
				}()
			}
			result := MustInt64(tt.input)
			if !tt.wantPanic && result != tt.expected {
				t.Errorf("MustInt64() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMustUint8(t *testing.T) {
	tests := []struct {
		name      string
		input     int
		expected  uint8
		wantPanic bool
	}{
		{"valid small value", 100, 100, false},
		{"valid zero", 0, 0, false},
		{"valid max uint8", math.MaxUint8, math.MaxUint8, false},
		{"negative value", -1, 0, true},
		{"negative large", -100, 0, true},
		{"overflow", math.MaxUint8 + 1, 0, true},
		{"large overflow", 1000, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("MustUint8() expected to panic but didn't")
					}
				}()
			}
			result := MustUint8(tt.input)
			if !tt.wantPanic && result != tt.expected {
				t.Errorf("MustUint8() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMustUint64(t *testing.T) {
	tests := []struct {
		name      string
		input     int64
		expected  uint64
		wantPanic bool
	}{
		{"valid small value", 100, 100, false},
		{"valid zero", 0, 0, false},
		{"valid max int64", math.MaxInt64, uint64(math.MaxInt64), false},
		{"negative value", -1, 0, true},
		{"negative large", -100, 0, true},
		{"min int64", math.MinInt64, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("MustUint64() expected to panic but didn't")
					}
				}()
			}
			result := MustUint64(tt.input)
			if !tt.wantPanic && result != tt.expected {
				t.Errorf("MustUint64() = %v, want %v", result, tt.expected)
			}
		})
	}
}
