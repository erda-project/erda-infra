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
	"fmt"
	"math"
)

// MustInt converts uint64 to int, panics if overflow
func MustInt(val uint64) int {
	if val > math.MaxInt {
		panic(fmt.Sprintf("numutil: uint64 value %d overflows int", val))
	}
	return int(val)
}

// MustInt64 converts uint64 to int64, panics if overflow
func MustInt64(val uint64) int64 {
	if val > math.MaxInt64 {
		panic(fmt.Sprintf("numutil: uint64 value %d overflows int64", val))
	}
	return int64(val)
}

// MustUint8 converts int to uint8, panics if overflow or negative
func MustUint8(val int) uint8 {
	if val < 0 {
		panic(fmt.Sprintf("numutil: int value %d is negative, cannot convert to uint8", val))
	}
	if val > math.MaxUint8 {
		panic(fmt.Sprintf("numutil: int value %d overflows uint8", val))
	}
	return uint8(val)
}

// MustUint64 converts int64 to uint64, panics if negative
func MustUint64(val int64) uint64 {
	if val < 0 {
		panic(fmt.Sprintf("numutil: int64 value %d is negative, cannot convert to uint64", val))
	}
	return uint64(val)
}
