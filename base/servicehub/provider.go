// Author: recallsong
// Email: songruiguo@qq.com

package servicehub

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-infra/base/logs"
)

// Creator .
type Creator func() Provider

// ProviderDefine .
type ProviderDefine interface {
	Service() []string
	Creator() Creator
}

// ProviderUsageSummary .
type ProviderUsageSummary interface {
	Summary() string
}

// ProviderUsage .
type ProviderUsage interface {
	Description() string
}

// ServiceDependencies .
type ServiceDependencies interface {
	Dependencies() []string
}

// ConfigCreator .
type ConfigCreator interface {
	Config() interface{}
}

// serviceProviders .
var serviceProviders = map[string]ProviderDefine{}

// RegisterProvider .
func RegisterProvider(name string, define ProviderDefine) {
	if _, ok := serviceProviders[name]; ok {
		fmt.Printf("provider %s already exist\n", name)
		os.Exit(-1)
	}
	serviceProviders[name] = define
}

// Provider .
type Provider interface{}

// Context .
type Context interface {
	Hub() *Hub
	Config() interface{}
	Logger() logs.Logger
	Service(name string, options ...interface{}) interface{}
}

// ProviderRunner .
type ProviderRunner interface {
	Start() error
	Close() error
}

// ProviderInitializer .
type ProviderInitializer interface {
	Init(ctx Context) error
}

// DependencyProvider .
type DependencyProvider interface {
	Provide(name string, options ...interface{}) interface{}
}
