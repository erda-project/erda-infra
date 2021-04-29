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

package servicehub

import "reflect"

// Spec define provider and register with Register function
type Spec struct {
	Define               interface{}           // optional
	Services             []string              // optional
	Dependencies         []string              // optional
	OptionalDependencies []string              // optional
	DependenciesFunc     func(h *Hub) []string // optional
	Summary              string                // optional
	Description          string                // optional
	ConfigFunc           func() interface{}    // optional
	Types                []reflect.Type        // optional
	Creator              Creator               // required
}

// Register .
func Register(name string, spec *Spec) {
	RegisterProvider(name, &specDefine{spec}) // wrap Spec as ProviderDefine
}

// ensure specDefine implements some interface
var (
	// _ ProviderDefine       = (*specDefine)(nil) // through RegisterProvider to ensure
	_ ProviderServices     = (*specDefine)(nil)
	_ ServiceTypes         = (*specDefine)(nil)
	_ ProviderUsageSummary = (*specDefine)(nil)
	_ ProviderUsage        = (*specDefine)(nil)
	_ ProviderUsage        = (*specDefine)(nil)
	_ ServiceDependencies  = (*specDefine)(nil)
	_ ConfigCreator        = (*specDefine)(nil)
	_ ConfigCreator        = (*specDefine)(nil)
)

type specDefine struct {
	s *Spec
}

func (d *specDefine) Services() []string {
	if len(d.s.Services) > 0 {
		return d.s.Services
	}
	if d, ok := d.s.Define.(ProviderServices); ok {
		return d.Services()
	}
	return nil
}

func (d *specDefine) Types() []reflect.Type {
	if len(d.s.Types) > 0 {
		return d.s.Types
	}
	if d, ok := d.s.Define.(ServiceTypes); ok {
		return d.Types()
	}
	return nil
}

func (d *specDefine) Dependencies(h *Hub) []string {
	var list = d.s.Dependencies
	for _, svr := range d.s.OptionalDependencies {
		if h.IsServiceExist(svr) {
			list = append(list, svr)
		}
	}
	if d.s.DependenciesFunc != nil {
		list = append(list, d.s.DependenciesFunc(h)...)
	}
	if len(list) > 0 {
		return list
	}
	if d, ok := d.s.Define.(ServiceDependencies); ok {
		return d.Dependencies(h)
	}
	return nil
}

func (d *specDefine) Summary() string {
	if len(d.s.Summary) > 0 {
		return d.s.Summary
	}
	if d, ok := d.s.Define.(ProviderUsageSummary); ok {
		return d.Summary()
	}
	return ""
}

func (d *specDefine) Description() string {
	if len(d.s.Description) > 0 {
		return d.s.Description
	}
	if d, ok := d.s.Define.(ProviderUsage); ok {
		return d.Description()
	}
	return ""
}

func (d *specDefine) Config() interface{} {
	if d.s.ConfigFunc != nil {
		return d.s.ConfigFunc()
	}
	if d, ok := d.s.Define.(ConfigCreator); ok {
		return d.Config()
	}
	return nil
}

func (d *specDefine) Creator() Creator {
	if d.s.Creator != nil {
		return d.s.Creator
	}
	if d, ok := d.s.Define.(ProviderDefine); ok {
		return d.Creator()
	}
	return nil // panic
}
