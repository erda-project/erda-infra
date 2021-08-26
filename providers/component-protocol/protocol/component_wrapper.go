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
	"context"
	"encoding/json"
	"reflect"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

func wrapCompRender(cr CompRender, ver string) CompRender {
	// 没版本的原始代码
	if ver == "" {
		return cr
	}
	return &compRenderWrapper{cr: cr}
}

type compRenderWrapper struct {
	cr CompRender
}

// Render .
func (w *compRenderWrapper) Render(ctx context.Context, c *cptype.Component, scenario cptype.Scenario, event cptype.ComponentEvent, gs *cptype.GlobalStateData) (err error) {
	if err = unmarshal(&w.cr, c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			// not marshal invoke fail
			return
		}
		err = marshal(&w.cr, c)
	}()
	err = w.cr.Render(ctx, c, scenario, event, gs)
	return
}

func unmarshal(cr *CompRender, c *cptype.Component) error {
	v, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return json.Unmarshal(v, cr)
}

func marshal(cr *CompRender, c *cptype.Component) error {
	var tmp cptype.Component
	v, err := json.Marshal(cr)
	if err != nil {
		return err
	}
	err = json.Unmarshal(v, &tmp)
	if err != nil {
		return err
	}
	tr := reflect.TypeOf(*cr).Elem()
	fields := getAllFields(c)
	for _, fieldName := range fields {
		if f, ok := tr.FieldByName(fieldName); ok {
			if tag := f.Tag.Get("json"); tag == "-" {
				continue
			}
			switch fieldName {
			case "Version":
				c.Version = tmp.Version
			case "Type":
				c.Type = tmp.Type
			case "Name":
				c.Name = tmp.Name
			case "Props":
				c.Props = tmp.Props
			case "Data":
				c.Data = tmp.Data
			case "State":
				c.State = tmp.State
			case "Operations":
				c.Operations = tmp.Operations
			}
		}
	}
	return nil
}

func getAllFields(o interface{}) (f []string) {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		f = append(f, t.Field(i).Name)
	}
	return
}
