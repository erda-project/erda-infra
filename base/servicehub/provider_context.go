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

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/recallsong/go-utils/config"
	"github.com/recallsong/go-utils/encoding/jsonx"
	"github.com/recallsong/unmarshal"
	unmarshalflag "github.com/recallsong/unmarshal/unmarshal-flag"
	"github.com/spf13/pflag"

	"github.com/erda-project/erda-infra/base/logs"
)

type inheritLabelStrategy string

const (
	inheritLabelTrue      inheritLabelStrategy = "true"
	inheritLabelFalse     inheritLabelStrategy = "false"
	inheritLabelPreferred inheritLabelStrategy = "preferred"
)

type providerContext struct {
	context.Context
	hub         *Hub
	key         string
	label       string
	name        string
	cfg         interface{}
	provider    Provider
	structValue reflect.Value
	structType  reflect.Type
	define      ProviderDefine
	tasks       []task
}

var loggerType = reflect.TypeOf((*logs.Logger)(nil)).Elem()

func (c *providerContext) BindConfig(flags *pflag.FlagSet) (err error) {
	if creator, ok := c.define.(ConfigCreator); ok {
		cfg := creator.Config()
		if cfg != nil {
			err = unmarshal.BindDefault(cfg)
			if err != nil {
				return err
			}
			if c.cfg != nil {
				err = config.ConvertData(c.cfg, cfg, "file")
				if err != nil {
					return err
				}
			}
			err = unmarshal.BindEnv(cfg)
			if err != nil {
				return err
			}
			if flags != nil {
				err = unmarshalflag.BindFlag(flags, cfg)
				if err != nil {
					return err
				}
			}
			c.cfg = cfg
			return nil
		}
	}
	c.cfg = nil
	return nil
}

func (c *providerContext) Init() (err error) {
	if reflect.ValueOf(c.provider).Kind() == reflect.Ptr && c.structType != nil {
		value, typ := c.structValue, c.structType
		var (
			cfgValue *reflect.Value
			cfgType  reflect.Type
		)
		if c.cfg != nil {
			value := reflect.ValueOf(c.cfg)
			cfgValue = &value
			cfgType = cfgValue.Type()
		}
		fields := typ.NumField()
		for i := 0; i < fields; i++ {
			if !value.Field(i).CanSet() {
				continue
			}
			field := typ.Field(i)
			if field.Type == loggerType {
				logger := c.Logger()
				value.Field(i).Set(reflect.ValueOf(logger))
			}
			if cfgValue != nil && field.Type == cfgType {
				value.Field(i).Set(*cfgValue)
			}
			service := field.Tag.Get("service")
			if len(service) <= 0 {
				service = field.Tag.Get("autowired")
			}
			if service == "-" {
				continue
			}
			service = c.adjustDependServiceLabel(service, &field)
			dc := newDependencyContext(
				service,
				c.name,
				field.Type,
				field.Tag,
			)
			instance := c.hub.getService(dc)
			if len(service) > 0 && instance == nil {
				opt, err := boolTagValue(field.Tag, "optional", false)
				if err != nil {
					return fmt.Errorf("invalid optional tag value in %s.%s: %s", typ.String(), field.Name, err)
				}
				if opt {
					continue
				}
				return fmt.Errorf("not found service %q", service)
			}
			if instance == nil {
				continue
			}
			if !reflect.TypeOf(instance).AssignableTo(field.Type) {
				return fmt.Errorf("service %q not implement %s", service, field.Type)
			}
			value.Field(i).Set(reflect.ValueOf(instance))
		}
	}
	if c.cfg != nil {
		key := c.key
		if key != c.name {
			key = fmt.Sprintf("%s (%s)", key, c.name)
		}
		if os.Getenv("LOG_LEVEL") == "debug" {
			fmt.Printf("provider %s config: \n%s\n", key, jsonx.MarshalAndIndent(c.cfg))
		}
		// c.hub.logger.Debugf("provider %s config: \n%s", key, jsonx.MarshalAndIndent(c.cfg))
	}

	if initializer, ok := c.provider.(ProviderInitializer); ok {
		err = initializer.Init(c)
		if err != nil {
			return fmt.Errorf("fail to Init provider %s: %s", c.name, err)
		}
	}
	return nil
}

// Define .
func (c *providerContext) Define() ProviderDefine {
	return c.define
}

func (c *providerContext) dependencies() string {
	services, providers := c.Dependencies()
	if len(services) > 0 && len(providers) > 0 {
		return fmt.Sprintf("services: %v, providers: %v", services, providers)
	} else if len(services) > 0 {
		return fmt.Sprintf("services: %v", services)
	} else if len(providers) > 0 {
		return fmt.Sprintf("providers: %v", providers)
	}
	return ""
}

func boolTagValue(tag reflect.StructTag, key string, defval bool) (bool, error) {
	opt, ok := tag.Lookup(key)
	if ok {
		if len(opt) > 0 {
			b, err := strconv.ParseBool(opt)
			if err != nil {
				return defval, err
			}
			return b, nil
		}
	}
	return defval, nil
}

func (c *providerContext) adjustDependServiceLabel(service string, field *reflect.StructField) string {
	if len(c.label) == 0 || strings.Contains(service, "@") {
		return service
	}
	inheritLabel := field.Tag.Get("inherit-label")
	switch inheritLabelStrategy(inheritLabel) {
	case inheritLabelTrue:
		return fmt.Sprintf("%s@%s", service, c.label)
	case inheritLabelPreferred:
		pcs := c.hub.servicesMap[service]
		for _, pc := range pcs {
			if pc.label == c.label {
				return fmt.Sprintf("%s@%s", service, c.label)
			}
		}
	case inheritLabelFalse:
	default:
	}
	return service
}

func (c *providerContext) fullName() string {
	if len(c.label) == 0 {
		return c.name
	}
	return fmt.Sprintf("%s@%s", c.name, c.label)
}

// Dependencies .
func (c *providerContext) Dependencies() (services []string, providers []string) {
	srvset, provset := make(map[string]bool), make(map[reflect.Type]bool)
	if deps, ok := c.define.(FixedServiceDependencies); ok {
		for _, service := range deps.Dependencies() {
			if !srvset[service] {
				services = append(services, service)
				srvset[service] = true
			}
		}
	}
	if deps, ok := c.define.(ServiceDependencies); ok {
		for _, service := range deps.Dependencies(c.hub) {
			if !srvset[service] {
				services = append(services, service)
				srvset[service] = true
			}
		}
	}
	if deps, ok := c.define.(OptionalServiceDependencies); ok {
		for _, service := range deps.OptionalDependencies(c.hub) {
			if len(c.hub.servicesMap[service]) > 0 && !srvset[service] {
				services = append(services, service)
				srvset[service] = true
			}
		}
	}
	if c.structType != nil {
		fields := c.structType.NumField()
		for i := 0; i < fields; i++ {
			field := c.structType.Field(i)
			service := field.Tag.Get("service")
			if len(service) <= 0 {
				service = field.Tag.Get("autowired")
			}
			if service == "-" {
				continue
			}
			if len(service) > 0 {
				service = c.adjustDependServiceLabel(service, &field)
				opt, _ := boolTagValue(field.Tag, "optional", false)
				if opt {
					if len(c.hub.servicesMap[service]) > 0 && !srvset[service] {
						services = append(services, service)
						srvset[service] = true
					}
				} else if !srvset[service] {
					services = append(services, service)
					srvset[service] = true
				}
				continue
			}
			if !c.structValue.Field(i).CanSet() {
				continue
			}
			plist := c.hub.servicesTypes[field.Type]
			if len(plist) > 0 && !provset[field.Type] {
				provset[field.Type] = true
				providers = append(providers, plist[0].name)
			}
		}
	}
	return
}

// Hub .
func (c *providerContext) Hub() *Hub {
	return c.hub
}

// Logger .
func (c *providerContext) Logger() logs.Logger {
	if c.hub.logger == nil {
		return nil
	}
	return c.hub.logger.Sub(c.key)
}

// Config .
func (c *providerContext) Config() interface{} {
	return c.cfg
}

// Service .
func (c *providerContext) Service(name string, options ...interface{}) interface{} {
	return c.hub.getService(newDependencyContext(
		name,
		c.name,
		nil,
		reflect.StructTag(""),
	), options...)
}

// AddTask .
func (c *providerContext) AddTask(fn func(context.Context) error, options ...TaskOption) {
	t := task{
		name: "",
		fn:   fn,
	}
	for _, opt := range options {
		opt(&t)
	}
	c.tasks = append(c.tasks, t)
}

// Label .
func (c *providerContext) Label() string {
	return c.label
}

// Key .
func (c *providerContext) Key() string {
	return c.key
}

// Provider .
func (c *providerContext) Provider() Provider {
	return c.provider
}

// WithTaskName .
func WithTaskName(name string) TaskOption {
	return func(t *task) {
		t.name = name
	}
}

type task struct {
	name string
	fn   func(context.Context) error
}

// dependencyContext .
type dependencyContext struct {
	typ     reflect.Type
	tags    reflect.StructTag
	service string
	key     string
	label   string
	caller  string
}

func (dc *dependencyContext) Type() reflect.Type      { return dc.typ }
func (dc *dependencyContext) Tags() reflect.StructTag { return dc.tags }
func (dc *dependencyContext) Service() string         { return dc.service }
func (dc *dependencyContext) Key() string             { return dc.key }
func (dc *dependencyContext) Label() string           { return dc.label }
func (dc *dependencyContext) Caller() string          { return dc.caller }

func newDependencyContext(service, caller string, typ reflect.Type, tags reflect.StructTag) *dependencyContext {
	dc := &dependencyContext{
		typ:     typ,
		tags:    tags,
		key:     service,
		service: service,
		caller:  caller,
	}
	idx := strings.Index(service, "@")
	if idx > 0 {
		dc.service = service[0:idx]
		dc.label = service[idx+1:]
	}
	return dc
}
