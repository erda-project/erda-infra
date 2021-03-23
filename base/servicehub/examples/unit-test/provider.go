// Author: recallsong
// Email: songruiguo@qq.com

package example

import (
	"github.com/erda-project/erda-infra/base/servicehub"
)

// Interface .
type Interface interface {
	Hello(name string) string
	Add(a, b int) int
}

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Service() []string { return []string{"example"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

var _ Interface = (*provider)(nil) // check interface implemented

type provider struct{}

func (p *provider) Hello(name string) string {
	return "hello " + name
}

func (p *provider) Add(a, b int) int {
	return a + b
}

func (p *provider) sub(a, b int) int {
	return a - b
}

func init() {
	servicehub.RegisterProvider("example-provider", &define{})
}
