// Copyright 2021 Terminus
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

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
)

// Creator .
type Creator func() Provider

// ProviderDefine .
type ProviderDefine interface {
	Creator() Creator
}

// ProviderService deprecated, use ProviderServices
type ProviderService interface {
	Service() []string
}

// ProviderServices .
type ProviderServices interface {
	Services() []string
}

// ServiceTypes .
type ServiceTypes interface {
	Types() []reflect.Type
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

// OptionalServiceDependencies .
type OptionalServiceDependencies interface {
	OptionalDependencies() []string
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

// ProviderRunnerWithContext .
type ProviderRunnerWithContext interface {
	Run(context.Context) error
}

// ProviderInitializer .
type ProviderInitializer interface {
	Init(ctx Context) error
}

// DependencyContext .
type DependencyContext interface {
	Type() reflect.Type
	Tags() reflect.StructTag
	Service() string
	Key() string
	Label() string
	Caller() string
}

// DependencyProvider .
type DependencyProvider interface {
	Provide(ctx DependencyContext, options ...interface{}) interface{}
}
