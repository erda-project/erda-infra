// Author: recallsong
// Email: songruiguo@qq.com

package servicehub

import "reflect"

// Spec define provider and register with Register function
type Spec struct {
	Services     []string           // optional
	Dependencies []string           // optional
	Summary      string             // optional
	Description  string             // optional
	ConfigFunc   func() interface{} // optional
	Types        []reflect.Type     // optional
	Creator      Creator            // required
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
	if d.s.ConfigFunc != nil {
		return d.s.ConfigFunc()
	}
	return nil
}

func (d *specDefine) Creator() Creator {
	return d.s.Creator
}
