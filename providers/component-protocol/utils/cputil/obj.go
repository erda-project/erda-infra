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

package cputil

import (
	"encoding/json"
	"fmt"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

//var json = jsi.ConfigFastest

// ObjJSONTransfer transfer from src to dst using json.
func ObjJSONTransfer(src interface{}, dst interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

// MustObjJSONTransfer .
func MustObjJSONTransfer(src interface{}, dst interface{}) interface{} {
	err := ObjJSONTransfer(src, dst)
	if err != nil {
		panic(fmt.Errorf("err: %v, src: %+v, dst: %+v", err, src, dst))
	}
	return dst
}

// MustConvertProps .
func MustConvertProps(props interface{}) cptype.ComponentProps {
	return *MustObjJSONTransfer(props, &cptype.ComponentProps{}).(*cptype.ComponentProps)
}

const (
	mapKeyMeta = "meta"
)

// MustFlatMapMeta .
func MustFlatMapMeta(input interface{}, removeMetaAfterFlat bool) {
	switch in := input.(type) {
	case map[string]interface{}:
		m := &in
		flatMapMeta(m, removeMetaAfterFlat)
	case *map[string]interface{}:
		m := in
		flatMapMeta(m, removeMetaAfterFlat)
	case []interface{}:
		for _, v := range in {
			v := v
			MustFlatMapMeta(v, removeMetaAfterFlat)
		}
	default:
		return
	}
}

func flatMapMeta(m *map[string]interface{}, removeMetaAfterFlat bool) {
	for k, v := range *m {
		v := v
		if k != mapKeyMeta {
			MustFlatMapMeta(v, removeMetaAfterFlat)
			(*m)[k] = v
			continue
		}
		// meta
		metaMap, ok := v.(map[string]interface{})
		if !ok {
			(*m)[k] = v
			continue
		}
		// flat map
		for k, v := range metaMap {
			v := v
			MustFlatMapMeta(v, removeMetaAfterFlat)
			(*m)[k] = v
		}
		if removeMetaAfterFlat {
			delete(*m, mapKeyMeta)
		}
	}
}
