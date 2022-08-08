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
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda-infra/base/logs"
)

type testBaseProvider struct{}

type testInitProvider struct {
	initialized chan interface{}
}

func (p *testInitProvider) Init(ctx Context) error {
	p.initialized <- nil
	return nil
}

type testRunProvider struct {
	running chan interface{}
	exited  chan interface{}
}

func (p *testRunProvider) Run(ctx context.Context) error {
	for {
		select {
		case p.running <- nil:
		case <-ctx.Done():
			p.exited <- nil
			return nil
		}
	}
}

type testStartProvider struct {
	started chan interface{}
	closed  chan interface{}
}

func (p *testStartProvider) Start() error {
	p.started <- nil
	return nil
}

func (p *testStartProvider) Close() error {
	p.closed <- nil
	return nil
}

type testDefine struct {
	name string
	spec *Spec
}

func testProviderName(name string) string { return "hub-" + name + "-provider" }

func testRegister(name string, deps []string, optdeps []string, creator func() Provider) testDefine {
	if creator == nil {
		creator = func() Provider {
			return struct{}{}
		}
	}
	return testDefine{
		testProviderName(name),
		&Spec{
			Services:             []string{name},
			Dependencies:         deps,
			OptionalDependencies: optdeps,
			Creator:              creator,
		},
	}
}

func testContent(names ...string) string {
	sb := strings.Builder{}
	for _, name := range names {
		sb.WriteString(testProviderName(name) + ":\n")
	}
	return sb.String()
}

func TestHub(t *testing.T) {
	type testItem struct {
		d         testDefine
		wait      func()
		check     func() error
		closeWait func()
	}
	type testFunc func() *testItem

	runProvider := func(name string) func() *testItem {
		return func() *testItem {
			createCh := make(chan interface{}, 1)
			configCh := make(chan interface{}, 1)
			initCh := make(chan interface{}, 1)
			startCh := make(chan interface{}, 1)
			closeCh := make(chan interface{}, 1)
			runCh := make(chan interface{}, 1)
			exitCh := make(chan interface{}, 1)
			p := &struct {
				testInitProvider
				testStartProvider
				testRunProvider
			}{
				testInitProvider:  testInitProvider{initCh},
				testStartProvider: testStartProvider{startCh, closeCh},
				testRunProvider:   testRunProvider{runCh, exitCh},
			}
			item := &testItem{
				d: testDefine{
					testProviderName(name),
					&Spec{
						ConfigFunc: func() interface{} {
							configCh <- nil
							return nil
						},
						Creator: func() Provider {
							createCh <- p
							return p
						},
					},
				},
				wait: func() {
					<-createCh
					<-configCh
					<-p.initialized
					<-p.started
					<-p.running
				},
				closeWait: func() {
					<-p.closed
					<-p.exited
				},
			}
			return item
		}
	}
	tests := []struct {
		name    string
		funcs   []testFunc
		content string
	}{
		{
			name: "empty",
			funcs: []testFunc{
				func() *testItem {
					createCh := make(chan interface{})
					item := &testItem{
						d: testDefine{
							testProviderName("test1"),
							&Spec{
								Creator: func() Provider {
									p := &testBaseProvider{}
									createCh <- p
									return p
								},
							},
						},
						wait: func() {
							<-createCh
						},
					}
					return item
				},
			},
			content: testContent("test1"),
		},
		{
			name: "config",
			funcs: []testFunc{
				func() *testItem {
					createCh := make(chan interface{})
					configCh := make(chan interface{})
					item := &testItem{
						d: testDefine{
							testProviderName("test1"),
							&Spec{
								ConfigFunc: func() interface{} {
									configCh <- nil
									return nil
								},
								Creator: func() Provider {
									p := &testBaseProvider{}
									createCh <- p
									return p
								},
							},
						},
						wait: func() {
							<-createCh
							<-configCh
						},
					}
					return item
				},
			},
			content: testContent("test1"),
		},
		{
			name: "init",
			funcs: []testFunc{
				func() *testItem {
					createCh := make(chan interface{})
					configCh := make(chan interface{})
					initCh := make(chan interface{})
					p := &testInitProvider{initCh}
					item := &testItem{
						d: testDefine{
							testProviderName("test1"),
							&Spec{
								ConfigFunc: func() interface{} {
									configCh <- nil
									return nil
								},
								Creator: func() Provider {
									createCh <- p
									return p
								},
							},
						},
						wait: func() {
							<-createCh
							<-configCh
							<-p.initialized
						},
					}
					return item
				},
			},
			content: testContent("test1"),
		},
		{
			name: "start",
			funcs: []testFunc{
				func() *testItem {
					createCh := make(chan interface{})
					configCh := make(chan interface{})
					initCh := make(chan interface{})
					startCh := make(chan interface{})
					closeCh := make(chan interface{})
					p := &struct {
						testInitProvider
						testStartProvider
					}{
						testInitProvider:  testInitProvider{initCh},
						testStartProvider: testStartProvider{startCh, closeCh},
					}
					item := &testItem{
						d: testDefine{
							testProviderName("test1"),
							&Spec{
								ConfigFunc: func() interface{} {
									configCh <- nil
									return nil
								},
								Creator: func() Provider {
									createCh <- p
									return p
								},
							},
						},
						wait: func() {
							<-createCh
							<-configCh
							<-p.initialized
							<-p.started
						},
						closeWait: func() {
							<-p.closed
						},
					}
					return item
				},
			},
			content: testContent("test1"),
		},
		{
			name: "run",
			funcs: []testFunc{
				runProvider("test1"),
			},
			content: testContent("test1"),
		},
		{
			name: "run many",
			funcs: []testFunc{
				runProvider("test1"),
				runProvider("test2"),
				runProvider("test3"),
			},
			content: testContent("test1", "test2", "test3"),
		},
		{
			name: "check config and log",
			funcs: []testFunc{
				func() *testItem {
					type config struct {
						Name string `default:"name1"`
					}
					type provider struct {
						Cfg *config
						Log logs.Logger
						testInitProvider
					}
					initCh := make(chan interface{})
					cfg := &config{}
					p := &provider{testInitProvider: testInitProvider{initCh}}
					item := &testItem{
						d: testDefine{
							testProviderName("test1"),
							&Spec{
								ConfigFunc: func() interface{} { return cfg },
								Creator:    func() Provider { return p },
							},
						},
						wait: func() {
							<-initCh
						},
						check: func() error {
							if p.Cfg != cfg {
								return fmt.Errorf("config field not match")
							}
							if p.Cfg.Name != "name1" {
								return fmt.Errorf("invalid config value, want config.name=%q, but got %q", "name1", p.Cfg.Name)
							}
							if p.Log == nil {
								return fmt.Errorf("log field not setup")
							}
							return nil
						},
					}
					return item
				},
			},
			content: testContent("test1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var list []*testItem
			for _, fn := range tt.funcs {
				item := fn()
				list = append(list, item)
				Register(item.d.name, item.d.spec)
			}
			hub := New()
			go func() {
				hub.RunWithOptions(&RunOptions{Content: tt.content})
			}()
			for _, item := range list {
				item.wait()
			}
			for _, item := range list {
				if item.check != nil {
					err := item.check()
					if err != nil {
						t.Errorf("check provider error: %v", err)
					}
				}
			}
			wg := &sync.WaitGroup{}
			for _, item := range list {
				if item.closeWait != nil {
					wg.Add(1)
					go func(item *testItem) {
						defer wg.Done()
						item.closeWait()
					}(item)
				}
			}
			if err := hub.Close(); err != nil {
				t.Errorf("Hub.Close() = %v, want nil", err)
			}
			wg.Wait()
			for _, item := range list {
				delete(serviceProviders, item.d.name)
			}
		})
	}
}

func TestHub_Dependencies(t *testing.T) {
	tests := []struct {
		name      string
		providers []testDefine
		content   string
		hasErr    bool
	}{
		{
			name: "Dependencies",
			providers: []testDefine{
				testRegister("test1", nil, nil, nil),
				testRegister("test2", []string{"test1"}, nil, nil),
			},
			content: testContent("test1", "test2"),
		},
		{
			name: "Miss Dependencies",
			providers: []testDefine{
				testRegister("test1", nil, nil, nil),
				testRegister("test2", []string{"test1"}, nil, nil),
			},
			content: testContent("test2"),
			hasErr:  true,
		},
		{
			name: "Dependencies And Optional Dependencies",
			providers: []testDefine{
				testRegister("test1", nil, nil, nil),
				testRegister("test2", []string{"test1"}, nil, nil),
				testRegister("test3", []string{"test1"}, []string{"test2", "test4"}, nil),
			},
			content: testContent("test1", "test2", "test3"),
		},
		{
			name: "Optional Dependencies",
			providers: []testDefine{
				testRegister("test1", nil, nil, nil),
				testRegister("test3", nil, []string{"test2", "test4"}, nil),
			},
			content: testContent("test1", "test3"),
		},
		{
			name: "Circular Dependency",
			providers: []testDefine{
				testRegister("test1", nil, nil, nil),
				testRegister("test2", nil, nil, nil),
				testRegister("test3", []string{"test2", "test4"}, nil, nil),
				testRegister("test4", []string{"test3"}, nil, nil),
			},
			content: testContent("test1", "testt2", "test3", "test4"),
			hasErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range tt.providers {
				Register(p.name, p.spec)
			}
			hub := New()
			events := hub.Events()
			go func() {
				hub.RunWithOptions(&RunOptions{Content: tt.content})
			}()
			err := <-events.Initialized()
			if (err != nil) != tt.hasErr {
				if tt.hasErr {
					t.Errorf("got error %q, want err != nil", err)
				} else {
					t.Errorf("got error %q, want err == nil", err)
				}
			}
			if err := hub.Close(); err != nil {
				t.Errorf("Hub.Close() = %v, want nil", err)
			}
			for _, p := range tt.providers {
				delete(serviceProviders, p.name)
			}
		})
	}
}

func Test_boolTagValue(t *testing.T) {
	type args struct {
		tag    reflect.StructTag
		key    string
		defval bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			args: args{
				tag:    reflect.StructTag(`test-key:""`),
				key:    "test-key",
				defval: false,
			},
			want: false,
		},
		{
			args: args{
				tag:    reflect.StructTag(`test-key:""`),
				key:    "test-key",
				defval: true,
			},
			want: true,
		},
		{
			args: args{
				tag:    reflect.StructTag(``),
				key:    "test-key",
				defval: true,
			},
			want: true,
		},
		{
			args: args{
				tag:    reflect.StructTag(`test-key:"true"`),
				key:    "test-key",
				defval: false,
			},
			want: true,
		},
		{
			args: args{
				tag:    reflect.StructTag(`test-key:"false"`),
				key:    "test-key",
				defval: true,
			},
			want: false,
		},
		{
			args: args{
				tag:    reflect.StructTag(`test-key:"error"`),
				key:    "test-key",
				defval: true,
			},
			want:    true,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := boolTagValue(tt.args.tag, tt.args.key, tt.args.defval)
			if (err != nil) != tt.wantErr {
				t.Errorf("boolTagValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("boolTagValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_addProblematicProvider(t *testing.T) {
	h := Hub{}
	assert.Equal(t, 0, len(h.problematicProviderNames))
	someProviders := []string{"p2", "p1", "p3"}
	var wg sync.WaitGroup
	for _, p := range someProviders {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			h.addProblematicProvider(p)
		}(p)
	}
	wg.Wait()
	assert.Equal(t, 3, len(h.problematicProviderNames))
	sort.Strings(h.problematicProviderNames)
	assert.Equal(t, []string{"p1", "p2", "p3"}, h.problematicProviderNames)
}
