// Author: recallsong
// Email: songruiguo@qq.com

package servicehub

import "reflect"

// Spec define provider and register with RegisterProviderSpec function
type Spec struct {
	Services     []string           // optional
	Dependencies []string           // optional
	Summary      string             // optional
	Description  string             // optional
	ConfigFunc   func() interface{} // optional
	Types        []reflect.Type     // optional
	Creator      Creator            // required
}

// RegisterProviderSpec .
func RegisterProviderSpec(name string, spec *Spec) {
	RegisterProvider(name, &specDefine{spec}) // wrap Spec as ProviderDefine
}

type specDefine struct {
	s *Spec
}

func (d *specDefine) Services() []string {
	return d.s.Services
}

func (d *specDefine) Types() []reflect.Type {
	return d.s.Types
}

func (d *specDefine) Dependencies() []string {
	return d.s.Dependencies
}

func (d *specDefine) Summary() string {
	return d.s.Summary
}

func (d *specDefine) Description() string {
	return d.s.Description
}

func (d *specDefine) Config() interface{} {
	return d.s.ConfigFunc()
}

func (d *specDefine) Creator() Creator {
	return d.s.Creator
}
