// Author: recallsong
// Email: songruiguo@qq.com

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
