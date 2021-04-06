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

package i18n

import (
	"fmt"
	"regexp"
)

// Template .
type Template struct {
	key     string
	content string
}

// NewTemplate .
func NewTemplate(key string, content string) *Template {
	return &Template{key: key, content: content}
}

// Render .
func (t *Template) Render(args ...interface{}) string {
	return fmt.Sprintf(t.content, args...)
}

// Key .
func (t *Template) Key() string { return t.key }

// Content .
func (t *Template) Content() string { return t.content }

// RenderByKey Render by key eg: {{keyName}}
func (t *Template) RenderByKey(params map[string]string) string {
	reg := regexp.MustCompile(`{{.+?}}`)
	result := reg.ReplaceAllStringFunc(t.content, func(s string) string {
		key := s[2 : len(s)-2]
		value, ok := params[key]
		if ok {
			return value
		}
		return s
	})
	return result
}
