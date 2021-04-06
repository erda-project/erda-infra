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
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/erda-project/erda-infra/base/logs"
)

var (
	initWait sync.WaitGroup
	step     []string

	runningWait sync.WaitGroup
	tasks       int32
	exits       int32

	cfg      = &testConfig{}
	provider = &test1Provider{}
)

type testConfig struct {
	Message string `file:"message" flag:"msg" default:"hi" desc:"message to show"`
}

type test1Provider struct {
	Log logs.Logger
	Cfg *testConfig
}

func (p *test1Provider) Init(ctx Context) error {
	defer initWait.Done()
	step = append(step, "init provider")
	return nil
}

func (p *test1Provider) Run(ctx context.Context) error {
	runningWait.Done()
	atomic.AddInt32(&tasks, 1)
	select {
	case <-ctx.Done():
		atomic.AddInt32(&exits, 1)
	}
	return nil
}

func (p *test1Provider) Start() error {
	runningWait.Done()
	atomic.AddInt32(&tasks, 1)
	return nil
}

func (p *test1Provider) Close() error {
	atomic.AddInt32(&exits, 1)
	return nil
}

func TestHub(t *testing.T) {
	runningWait.Add(2)
	initWait.Add(3)
	Register("hub-test-provider", &Spec{
		Services:    []string{"test"},
		Description: "this is provider for test",
		ConfigFunc: func() interface{} {
			defer initWait.Done()
			step = append(step, "create config")
			return cfg
		},
		Creator: func() Provider {
			defer initWait.Done()
			step = append(step, "create provider")
			return provider
		},
	})

	var wg sync.WaitGroup
	wg.Add(1)
	hub := New()
	go func() {
		defer wg.Done()
		hub.RunWithOptions(&RunOptions{
			Content: `
hub-test-provider:
    message: "hello world"
`})
	}()
	initWait.Wait()
	if strings.Join(step, ",") != "create provider,create config,init provider" {
		t.Errorf("out-of-order init step, got %q", strings.Join(step, ","))
	}
	runningWait.Wait()
	hub.Close()

	if provider.Log == nil {
		t.Errorf("logger is nil")
	}
	if provider.Cfg != cfg {
		t.Errorf("config got %v, but want %v", provider.Cfg, cfg)
	}
	if cfg.Message != "hello world" {
		t.Errorf("read config error, got cfg.Message = %q, but want %q", cfg.Message, "hello world")
	}

	wg.Wait()
	if tasks != 2 {
		t.Errorf("tasks(%d) != %d, some function not called", tasks, 2)
	}
	if tasks != exits {
		t.Errorf("tasks(%d) != exist(%d), maybe some function not exit", tasks, exits)
	}
}

func TestHub_Dependencies(t *testing.T) {
	Register("hub-test1-deps-provider", &Spec{
		Services: []string{"test1"},
		Creator: func() Provider {
			return struct{}{}
		},
	})
	Register("hub-test2-deps-provider", &Spec{
		Services:     []string{"test2"},
		Dependencies: []string{"test1"},
		Creator: func() Provider {
			return struct{}{}
		},
	})
	Register("hub-test3-deps-provider", &Spec{
		Services:             []string{"test3"},
		Dependencies:         []string{"test1"},
		OptionalDependencies: []string{"test2", "test4"},
		Creator: func() Provider {
			return struct{}{}
		},
	})
	Register("hub-test4-deps-provider", &Spec{
		Services:     []string{"test4"},
		Dependencies: []string{"test3"},
		Creator: func() Provider {
			return struct{}{}
		},
	})

	tests := []struct {
		name    string
		content string
		hasErr  bool
	}{
		{
			name: "Dependencies",
			content: `
hub-test1-deps-provider:
hub-test2-deps-provider:
`,
		},
		{
			name: "Miss Dependencies",
			content: `
hub-test2-deps-provider:
`,
			hasErr: true,
		},
		{
			name: "Dependencies And Optional Dependencies",
			content: `
hub-test1-deps-provider:
hub-test2-deps-provider:
hub-test3-deps-provider:
`,
		},
		{
			name: "Optional Dependencies",
			content: `
hub-test1-deps-provider:
hub-test3-deps-provider:
`,
		},
		{
			name: "Circular Dependency",
			content: `
hub-test1-deps-provider:
hub-test3-deps-provider:
hub-test4-deps-provider:
`,
			hasErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errCh := make(chan error)
			hub := New(WithListener(&DefaultListener{
				BeforeExitFunc: func(h *Hub, err error) error {
					errCh <- err
					return nil
				},
			}))
			go func() {
				hub.RunWithOptions(&RunOptions{Content: tt.content})
			}()
			err := <-errCh
			if (err != nil) != tt.hasErr {
				if tt.hasErr {
					t.Errorf("got error %v, want err != nil", err)
				} else {
					t.Errorf("got error %v, want err == nil", err)
				}
			}
		})
	}
}
