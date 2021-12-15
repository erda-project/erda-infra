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

package cptype

import (
	"strconv"

	"github.com/erda-project/erda-infra/pkg/strutil"
)

// Extra is a magic fields, and will be flat if specified in protocol.
type Extra struct {
	Extra ExtraMap `json:"extra,omitempty"`
}

// ExtraMap .
type ExtraMap map[string]interface{}

// Uint64 .
func (e Extra) Uint64(key string) uint64 {
	return e.Extra.Uint64(key)
}

// String .
func (e Extra) String(key string) string {
	return e.Extra.String(key)
}

// Uint64 .
func (em ExtraMap) Uint64(key string) uint64 {
	if len(em) == 0 {
		return 0
	}
	v, ok := em[key]
	if !ok {
		return 0
	}
	switch vv := v.(type) {
	case float64:
		return uint64(vv)
	case string:
		r, _ := strconv.ParseUint(vv, 10, 64)
		return r
	}
	return 0
}

// String .
func (em ExtraMap) String(key string) string {
	if len(em) == 0 {
		return ""
	}
	v, ok := em[key]
	if !ok {
		return ""
	}
	return strutil.String(v)
}

// Get .
func (em ExtraMap) Get(key string) interface{} {
	if len(em) == 0 {
		return nil
	}
	v, ok := em[key]
	if !ok {
		return nil
	}
	return v
}
