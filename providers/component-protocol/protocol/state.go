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

package protocol

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// setGlobalStateKV .
func setGlobalStateKV(p *cptype.ComponentProtocol, key string, value interface{}) {
	if p.GlobalState == nil {
		var gs = make(cptype.GlobalStateData)
		p.GlobalState = &gs
	}
	s := p.GlobalState
	(*s)[key] = value
}

// GetGlobalStateKV .
func GetGlobalStateKV(p *cptype.ComponentProtocol, key string) interface{} {
	if p.GlobalState == nil {
		return nil
	}
	return (*p.GlobalState)[key]
}

// getCompStateKV .
func getCompStateKV(c *cptype.Component, stateKey string) (interface{}, error) {
	if c == nil {
		err := fmt.Errorf("empty component")
		return nil, err
	}
	if _, ok := c.State[stateKey]; !ok {
		err := fmt.Errorf("state key [%s] not exist in component [%s] state", stateKey, c.Name)
		return nil, err
	}
	return c.State[stateKey], nil
}

// setCompStateValueFromComps .
func setCompStateValueFromComps(c *cptype.Component, key string, value interface{}) error {
	if c == nil {
		err := fmt.Errorf("empty component")
		return err
	}
	if key == "" {
		err := fmt.Errorf("empty state key")
		return err
	}
	if v, ok := c.State[key]; ok {
		logrus.Infof("state key already exist in component, component:%s, key:%s, value old:%+v, new:%+v", c.Name, key, v, value)
	}
	if c.State == nil {
		c.State = map[string]interface{}{}
	}
	c.State[key] = value
	return nil
}

// parseStateBound .
func parseStateBound(b string) (comp, key string, err error) {
	prefix := "{{"
	suffix := "}}"
	if !strings.HasPrefix(b, prefix) {
		err = fmt.Errorf("state bound not prefix with {{")
		return
	}
	if !strings.HasSuffix(b, "}}") {
		err = fmt.Errorf("state bound not suffix with }}")
		return
	}
	b = strings.TrimPrefix(b, prefix)
	b = strings.TrimPrefix(b, " ")
	b = strings.TrimSuffix(b, suffix)
	b = strings.TrimSuffix(b, " ")
	s := strings.Split(b, ".")
	if len(s) != 2 {
		err = fmt.Errorf("invalide bound expression: %s, with not exactly one '.' ", b)
		return
	}
	comp = s[0]
	key = s[1]
	return
}
