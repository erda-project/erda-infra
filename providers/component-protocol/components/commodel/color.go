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

package commodel

import (
	"fmt"
)

// UnifiedColor .
type UnifiedColor int

// IUnifiedColor .
type IUnifiedColor interface {
	fmt.Stringer
}

// UnifiedColor .
const (
	ColorRed UnifiedColor = iota
	ColorYellow
	ColorGreen
	ColorBlue
	ColorPurple
	ColorDefault
)

// String .
func (c UnifiedColor) String() string {
	switch c {
	case 0:
		return "red"
	case 1:
		return "yellow"
	case 2:
		return "green"
	case 3:
		return "blue"
	case 4:
		return "purple"
	default:
		return "default"
	}
}
