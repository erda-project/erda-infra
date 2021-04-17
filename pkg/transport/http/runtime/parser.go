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

package runtime

import (
	"fmt"
	"strings"

	"github.com/erda-project/erda-infra/pkg/transport/http/httprule"
)

// Matcher .
type Matcher interface {
	Match(path string) (map[string]string, error)
	IsStatic() bool
	Fields() []string
	Pattern() string
}

// Compile .
func Compile(path string) (Matcher, error) {
	if len(path) <= 1 {
		return &staticMacher{path}, nil
	}
	compiler, err := httprule.Parse(path)
	if err != nil {
		return nil, ErrInvalidPattern
	}
	temp := compiler.Compile()
	if len(temp.Fields) <= 0 {
		return &staticMacher{path}, nil
	}
	pattern, err := NewPattern(httprule.SupportPackageIsVersion1, temp.OpCodes, temp.Pool, temp.Verb)
	if err != nil {
		return nil, fmt.Errorf("fail to create path pattern: %s", err)
	}
	return &paramsMatcher{&pattern, path, temp.Fields}, nil
}

type staticMacher struct {
	path string
}

func (m *staticMacher) Match(path string) (map[string]string, error) {
	if m.path == path {
		return nil, nil
	}
	return nil, ErrNotMatch
}

func (m *staticMacher) IsStatic() bool   { return true }
func (m *staticMacher) Fields() []string { return nil }
func (m *staticMacher) Pattern() string  { return m.path }

type paramsMatcher struct {
	p      *Pattern
	path   string
	fields []string
}

func (m *paramsMatcher) Match(path string) (map[string]string, error) {
	if len(path) > 0 {
		components := strings.Split(path[1:], "/")
		last := len(components) - 1
		var verb string
		if idx := strings.LastIndex(components[last], ":"); idx >= 0 {
			c := components[last]
			components[last], verb = c[:idx], c[idx+1:]
		}
		return m.p.Match(components, verb)
	}
	return nil, ErrNotMatch
}

func (m *paramsMatcher) IsStatic() bool   { return false }
func (m *paramsMatcher) Fields() []string { return m.fields }
func (m *paramsMatcher) Pattern() string  { return m.path }
