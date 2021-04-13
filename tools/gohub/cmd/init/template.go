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

package init

type tempContext struct {
	Package     string
	Provider    string
	Description string
}

type providerTemplate struct {
	Content     string
	TestContext string
}

var templates = map[string]providerTemplate{
	"full": {
		Content:     simpleTemplate,     // TODO: full template
		TestContext: simpleTestTemplate, // TODO: full template
	},
	"simple": {
		Content:     simpleTemplate,
		TestContext: simpleTestTemplate,
	},
}

const (
	fullTemplate = ``
)

const (
	simpleTemplate = `package {{.Package}}

import (
	"context"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

type config struct {
	// some fields of config for this provider
	Message string ` + "`file:\"message\" flag:\"msg\" default:\"hi\" desc:\"message to print\"`" + ` 
}

// +provider
type provider struct {
	Cfg *config
	Log logs.Logger
}

// Run this is an optional
func (p *provider) Init(ctx servicehub.Context) error {
	p.Log.Info("message: ", p.Cfg.Message)
	return nil
}

// Run this is an optional
func (p *provider) Run(ctx context.Context) error {
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			p.Log.Info("do something...")
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.Register({{quote .Provider}}, &servicehub.Spec{
		Services:    []string{
			{{printf "%s-service" .Provider | quote}},
		},
		Description: {{printf "here is description of %s" .Provider | quote}},
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
`
	simpleTestTemplate = `package {{.Package}}

import (
	"fmt"
	"testing"

	"github.com/erda-project/erda-infra/base/servicehub"
)
	
type testInterface interface {
	testFunc(arg interface{}) interface{}
}
	
func (p *provider) testFunc(arg interface{}) interface{} {
	return fmt.Sprintf("%s -> result", arg)
}
	
func Test_provider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		config   string
		arg      interface{}
		want     interface{}
	}{
		{
			"case 1",
			{{quote .Provider}},
			` + "`\n" +
		`{{.Provider}}:
    message: "hello"
` +
		"`,\n" + `
			"test arg",
			"test arg -> result",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hub := servicehub.New()
			events := hub.Events()
			go func() {
				hub.RunWithOptions(&servicehub.RunOptions{Content: tt.config})
			}()
			<-events.Started()
	
			p := hub.Provider(tt.provider).(*provider)
			if got := p.testFunc(tt.arg); got != tt.want {
				t.Errorf("provider.testFunc() = %v, want %v", got, tt.want)
			}
			if err := hub.Close(); err != nil {
				t.Errorf("Hub.Close() = %v, want nil", err)
			}
		})
	}
}
	
func Test_provider_service(t *testing.T) {
	tests := []struct {
		name    string
		service string
		config  string
		arg     interface{}
		want    interface{}
	}{
		{
			"case 1",
			{{printf "%s-service" .Provider | quote}},
			` + "`\n" +
		`{{.Provider}}:
    message: "hello"
` +
		"`,\n" + `
			"test arg",
			"test arg -> result",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hub := servicehub.New()
			events := hub.Events()
			go func() {
				hub.RunWithOptions(&servicehub.RunOptions{Content: tt.config})
			}()
			<-events.Started()
			s := hub.Service(tt.service).(testInterface)
			if got := s.testFunc(tt.arg); got != tt.want {
				t.Errorf("(service %q).testFunc() = %v, want %v", tt.service, got, tt.want)
			}
			if err := hub.Close(); err != nil {
				t.Errorf("Hub.Close() = %v, want nil", err)
			}
		})
	}
}
`
)
