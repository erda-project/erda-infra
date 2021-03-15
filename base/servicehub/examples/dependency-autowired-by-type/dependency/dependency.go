// Author: recallsong
// Email: songruiguo@qq.com

package dependency

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

// Interface .
type Interface interface {
	Hello(name string) string
}

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Service() []string { return []string{"example-dependency"} }

// Declare what service types the provider provides
func (d *define) Types() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Interface)(nil)).Elem()}
}

// Describe information about this provider
func (d *define) Description() string { return "dependency for example" }

// Return an instance representing the configuration
func (d *define) Config() interface{} { return &config{} }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type config struct{}

type provider struct {
	C *config
	L logs.Logger
}

func (p *provider) Hello(name string) string {
	return "hello " + name
}

func init() {
	servicehub.RegisterProvider("example-dependency-provider", &define{})
}
