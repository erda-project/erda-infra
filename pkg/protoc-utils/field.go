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

package protocutils

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// GetFieldPath .
func GetFieldPath(key string, fields []*protogen.Field) (string, *protogen.Field, error) {
	keys := strings.Split(key, ".")
	names := make([]string, len(keys))
	last := len(keys) - 1
	var pfield *protogen.Field

	var fn func(keys []string, fields []*protogen.Field) error
	fn = func(keys []string, fields []*protogen.Field) error {
		if len(keys) <= 0 {
			return nil
		}
		for _, field := range fields {
			if string(field.Desc.Name()) == keys[0] {
				if len(keys[1:]) > 0 {
					if field.Message == nil {
						break
					}
					return fn(keys[1:], field.Message.Fields)
				}
				names[last] = field.GoName
				if pfield == nil {
					pfield = field
				}
				last--
				return nil
			}
		}
		return fmt.Errorf("field %q not exist", key)
	}
	err := fn(strings.Split(key, "."), fields)
	if err != nil {
		return "", nil, err
	}
	return strings.Join(names, "."), pfield, nil
}
