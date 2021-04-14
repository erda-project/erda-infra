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
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/logs/logrusx"
	graph "github.com/erda-project/erda-infra/base/servicehub/dependency-graph"
	"github.com/recallsong/go-utils/config"
	"github.com/recallsong/go-utils/encoding/jsonx"
	"github.com/recallsong/go-utils/errorx"
	"github.com/recallsong/go-utils/os/signalx"
	"github.com/recallsong/unmarshal"
	unmarshalflag "github.com/recallsong/unmarshal/unmarshal-flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Hub .
type Hub struct {
	logger        logs.Logger
	providersMap  map[string][]*providerContext
	providers     []*providerContext
	servicesMap   map[string][]*providerContext
	servicesTypes map[reflect.Type][]*providerContext
	lock          sync.RWMutex

	started bool
	cancel  func()
	wg      sync.WaitGroup

	listeners []Listener
}

// New .
func New(options ...interface{}) *Hub {
	hub := &Hub{}
	for _, opt := range options {
		processOptions(hub, opt)
	}
	if hub.logger == nil {
		level := os.Getenv("LOG_LEVEL")
		lvl, err := logrus.ParseLevel(level)
		if err == nil {
			hub.logger = logrusx.New(logrusx.WithLevel(lvl))
		} else {
			hub.logger = logrusx.New()
		}
	}
	return hub
}

// Init .
func (h *Hub) Init(config map[string]interface{}, flags *pflag.FlagSet, args []string) (err error) {
	defer func() {
		// exp := recover()
		// if exp != nil {
		// 	if e, ok := exp.(error); ok {
		// 		err = e
		// 	} else {
		// 		err = fmt.Errorf("%v", exp)
		// 	}
		// }
		if err != nil {
			h.logger.Errorf("fail to init service hub: %s", err)
		}
	}()
	for i, l := 0, len(h.listeners); i < l; i++ {
		err = h.listeners[i].BeforeInitialization(h, config)
		if err != nil {
			return err
		}
	}
	err = h.loadProviders(config)
	if err != nil {
		return err
	}

	depGraph, err := h.resolveDependency(h.providersMap)
	if err != nil {
		return fmt.Errorf("fail to resolve dependency: %s", err)
	}

	flags.BoolP("providers", "p", false, "print all providers supported")
	flags.BoolP("graph", "g", false, "print providers dependency graph")
	for _, ctx := range h.providers {
		err = ctx.BindConfig(flags)
		if err != nil {
			return fmt.Errorf("fail to bind config for provider %s: %s", ctx.name, err)
		}
	}
	err = flags.Parse(args)
	if err != nil {
		return fmt.Errorf("fail to bind flags: %s", err)
	}
	if ok, err := flags.GetBool("providers"); err == nil && ok {
		usage := Usage()
		fmt.Println(usage)
		os.Exit(0)
	}
	if ok, err := flags.GetBool("graph"); err == nil && ok {
		depGraph.Display()
		os.Exit(0)
	}
	for _, ctx := range h.providers {
		err = ctx.Init()
		if err != nil {
			return err
		}
		dependencies := ctx.dependencies()
		if len(dependencies) > 0 {
			h.logger.Infof("provider %s (depends %s) initialized", ctx.name, dependencies)
		} else {
			h.logger.Infof("provider %s initialized", ctx.name)
		}
	}
	for i := len(h.listeners) - 1; i >= 0; i-- {
		err = h.listeners[i].AfterInitialization(h)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Hub) resolveDependency(providersMap map[string][]*providerContext) (graph.Graph, error) {
	services := map[string][]*providerContext{}
	types := map[reflect.Type][]*providerContext{}
	for _, p := range providersMap {
		d := p[0].define
		var list []string
		if ps, ok := d.(ProviderServices); ok {
			list = ps.Services()
		} else if ps, ok := d.(ProviderService); ok {
			list = ps.Service()
		}
		for _, s := range list {
			if exist, ok := services[s]; ok {
				return nil, fmt.Errorf("service %s conflict between %s and %s", s, exist[0].name, p[0].name)
			}
			services[s] = p
		}
		if ts, ok := d.(ServiceTypes); ok {
			for _, t := range ts.Types() {
				if exist, ok := types[t]; ok {
					return nil, fmt.Errorf("service type %s conflict between %s and %s", t, exist[0].name, p[0].name)
				}
				types[t] = p
			}
		}
	}
	h.servicesMap = services
	h.servicesTypes = types
	var depGraph graph.Graph
	for name, p := range providersMap {
		providers := map[string]*providerContext{}
		dependsServices, dependsProviders := p[0].Dependencies()
	loop:
		for _, service := range dependsServices {
			name := service
			var label string
			idx := strings.Index(service, "@")
			if idx > 0 {
				name, label = service[0:idx], service[idx+1:]
			}
			if deps, ok := services[name]; ok {
				if len(label) > 0 {
					for _, dep := range deps {
						if dep.label == label {
							providers[dep.name] = dep
							continue loop
						}
					}
				} else if len(deps) > 0 {
					providers[deps[0].name] = deps[0]
					continue loop
				}
			}
			return nil, fmt.Errorf("miss provider of service %s", service)
		}
		node := graph.NewNode(name)
		for dep := range providers {
			node.Deps = append(node.Deps, dep)
		}
		for _, dep := range dependsProviders {
			if _, ok := providers[dep]; !ok {
				node.Deps = append(node.Deps, dep)
			}
		}
		depGraph = append(depGraph, node)
	}
	resolved, err := graph.Resolve(depGraph)
	if err != nil {
		depGraph.Display()
		return depGraph, err
	}
	var providers []*providerContext
	for _, node := range resolved {
		providers = append(providers, providersMap[node.Name]...)
	}
	h.providers = providers
	return resolved, nil
}

// StartWithSignal .
func (h *Hub) StartWithSignal() error {
	sigs := []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
	h.logger.Infof("signals to quit: %v", sigs)
	return h.Start(signalx.Notify(sigs...))
}

// Start .
func (h *Hub) Start(closer ...<-chan os.Signal) (err error) {
	h.lock.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	h.cancel = cancel
	ch := make(chan error, len(h.providers))
	var num int
	for _, item := range h.providers {
		if runner, ok := item.provider.(ProviderRunner); ok {
			num++
			h.wg.Add(1)
			go func(key, name string, provider ProviderRunner) {
				if key != name {
					key = fmt.Sprintf("%s (%s)", key, name)
				}
				h.logger.Infof("provider %s starting ...", key)
				err := provider.Start()
				if err != nil {
					h.logger.Errorf("fail to start provider %s: %s", key, err)
				} else {
					h.logger.Infof("provider %s closed", key)
				}
				h.wg.Done()
				ch <- err
			}(item.key, item.name, runner)
		}
		if runner, ok := item.provider.(ProviderRunnerWithContext); ok {
			num++
			h.wg.Add(1)
			go func(key, name string, provider ProviderRunnerWithContext) {
				if key != name {
					key = fmt.Sprintf("%s (%s)", key, name)
				}
				h.logger.Infof("provider %s running ...", key)
				err := provider.Run(ctx)
				if err != nil {
					h.logger.Errorf("fail to run provider %s: %s", key, err)
				} else {
					h.logger.Infof("provider %s exit", key)
				}
				h.wg.Done()
				ch <- err
			}(item.key, item.name, runner)
		}
	}
	h.started = true
	h.lock.Unlock()
	runtime.Gosched()

	for i, l := 0, len(h.listeners); i < l; i++ {
		err = h.listeners[i].AfterStart(h)
		if err != nil {
			return err
		}
	}

	closeCh, closed := make(chan struct{}), false
	var elock sync.Mutex
	for _, ch := range closer {
		go func(ch <-chan os.Signal) {
			select {
			case <-ch:
			case <-closeCh:
			}
			elock.Lock()
			fmt.Println()
			wait := make(chan error)
			go func() {
				wait <- h.Close()
			}()
			select {
			case <-time.After(30 * time.Second):
				h.logger.Errorf("exit service manager timeout !")
				os.Exit(1)
			case err := <-wait:
				if err != nil {
					h.logger.Errorf("fail to exit: %s", err)
					os.Exit(1)
				}
			}
		}(ch)
	}
	// wait to stop
	errs := errorx.Errors{}
	for i := 0; i < num; i++ {
		select {
		case err := <-ch:
			if err != nil {
				errs = append(errs, err)
				if !closed {
					close(closeCh)
					closed = true
				}
			}
		}
	}
	err = errs.MaybeUnwrap()
	for i, l := 0, len(h.listeners); i < l; i++ {
		err = h.listeners[i].BeforeExit(h, err)
	}
	return err
}

// Close .
func (h *Hub) Close() error {
	h.lock.Lock()
	if !h.started {
		h.lock.Unlock()
		return nil
	}
	var errs errorx.Errors
	for i := len(h.providers) - 1; i >= 0; i-- {
		if runner, ok := h.providers[i].provider.(ProviderRunner); ok {
			err := runner.Close()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	h.cancel()
	h.wg.Wait()
	h.started = false
	h.lock.Unlock()
	return errs.MaybeUnwrap()
}

type providerContext struct {
	hub         *Hub
	key         string
	label       string
	name        string
	cfg         interface{}
	provider    Provider
	structValue reflect.Value
	structType  reflect.Type
	define      ProviderDefine
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
			err = unmarshalflag.BindFlag(flags, cfg)
			if err != nil {
				return err
			}
			c.cfg = cfg
		}
	}
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
			dc := newDependencyContext(
				service,
				c.name,
				field.Type,
				field.Tag,
			)
			instance := c.hub.getService(dc)
			if len(service) > 0 && instance == nil {
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

// Dependencies .
func (c *providerContext) Dependencies() (services []string, providers []string) {
	srvset, provset := make(map[string]bool), make(map[reflect.Type]bool)
	if deps, ok := c.define.(ServiceDependencies); ok {
		for _, service := range deps.Dependencies() {
			if !srvset[service] {
				services = append(services, service)
				srvset[service] = true
			}
		}
	}
	if deps, ok := c.define.(OptionalServiceDependencies); ok {
		for _, service := range deps.OptionalDependencies() {
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
				if !srvset[service] {
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
	return c.hub.logger.Sub(c.name)
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

// Service .
func (h *Hub) Service(name string, options ...interface{}) interface{} {
	return h.getService(newDependencyContext(
		name,
		"",
		nil,
		reflect.StructTag(""),
	), options...)
}

func (h *Hub) getService(dc DependencyContext, options ...interface{}) (instance interface{}) {
	var pc *providerContext
	if len(dc.Service()) > 0 {
		if providers, ok := h.servicesMap[dc.Service()]; ok {
			if len(providers) > 0 {
				if len(dc.Label()) > 0 {
					for _, item := range providers {
						if item.label == dc.Label() {
							pc = item
							break
						}
					}
				} else {
					for _, item := range providers {
						if item.key == item.name {
							pc = item
							break
						}
					}
					if pc == nil && len(providers) > 0 {
						pc = providers[0]
					}
				}
			}
		}
	} else if dc.Type() != nil {
		providers := h.servicesTypes[dc.Type()]
		for _, item := range providers {
			if item.key == item.name {
				pc = item
				break
			}
		}
		if pc == nil && len(providers) > 0 {
			pc = providers[0]
		}
	}
	if pc != nil {
		provider := pc.provider
		if prod, ok := provider.(DependencyProvider); ok {
			return prod.Provide(dc, options...)
		}
		return provider
	}
	return nil
}

// Provider .
func (h *Hub) Provider(name string) interface{} {
	var label string
	idx := strings.Index(name, "@")
	if idx > 0 {
		label = name[idx+1:]
		name = name[0:idx]
	}
	ps := h.providersMap[name]
	if len(label) > 0 {
		for _, p := range ps {
			if p.label == label {
				return p.provider
			}
		}
	} else if len(ps) > 0 {
		return ps[0].provider
	}
	return nil
}

// RunOptions .
type RunOptions struct {
	Name       string
	ConfigFile string
	Content    interface{}
	Format     string
	Args       []string
}

// RunWithOptions .
func (h *Hub) RunWithOptions(opts *RunOptions) {
	name := opts.Name
	if len(name) <= 0 {
		name = getAppName(opts.Args...)
	}
	config.LoadEnvFile()

	var err error
	var start bool
	defer func() {
		if !start {
			for i, l := 0, len(h.listeners); i < l; i++ {
				err = h.listeners[i].BeforeExit(h, err)
			}
		}
		if err != nil {
			os.Exit(1)
		}
	}()

	format := "yaml"
	if len(opts.Format) > 0 {
		format = opts.Format
	}
	cfgmap := make(map[string]interface{})
	if opts.Content != nil {
		var reader io.Reader
		switch val := opts.Content.(type) {
		case map[string]interface{}:
			cfgmap = val
		case string:
			reader = strings.NewReader(val)
		case []byte:
			reader = bytes.NewReader(val)
		default:
			err = fmt.Errorf("invalid config content type")
			h.logger.Error(err)
			return
		}
		if reader != nil {
			err = config.UnmarshalToMap(reader, format, cfgmap)
			if err != nil {
				h.logger.Errorf("fail to parse %s config: %s", format, err)
				return
			}
		}
	}

	cfgfile := opts.ConfigFile
	if len(cfgmap) <= 0 && len(opts.ConfigFile) <= 0 {
		cfgfile = name + "." + format
	}
	cfgmap, err = h.loadConfigWithArgs(cfgfile, cfgmap, opts.Args...)
	if err != nil {
		return
	}

	flags := pflag.NewFlagSet(name, pflag.ExitOnError)
	flags.StringP("config", "c", cfgfile, "config file to load providers")
	err = h.Init(cfgmap, flags, opts.Args)
	if err != nil {
		return
	}
	defer h.Close()
	start = true
	err = h.StartWithSignal()
	if err != nil {
		return
	}
}

// Run .
func (h *Hub) Run(name, cfgfile string, args ...string) {
	h.RunWithOptions(&RunOptions{
		Name:       name,
		ConfigFile: cfgfile,
		Args:       args,
	})
}

// Run .
func Run(opts *RunOptions) *Hub {
	hub := New()
	hub.RunWithOptions(opts)
	return hub
}

func getAppName(args ...string) string {
	if len(args) <= 0 {
		return ""
	}
	name := args[0]
	idx := strings.LastIndex(os.Args[0], "/")
	if idx >= 0 {
		return name[idx+1:]
	}
	return ""
}
