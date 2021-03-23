// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Services() []string { return []string{"hello"} }

// Describe information about this provider
func (d *define) Description() string { return "hello for example" }

// Return an instance representing the configuration
func (d *define) Config() interface{} { return &config{} }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type config struct {
	Error bool `file:"error"`
}

type provider struct {
	Cfg *config
	Log logs.Logger
}

func (p *provider) Run(ctx context.Context) error {
	if p.Cfg.Error {
		time.Sleep(3 * time.Second)
		return fmt.Errorf("run error")
	}
	p.Log.Info("run with no error")
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.RegisterProvider("hello-provider", &define{})
}

func main() {
	servicehub.Run(&servicehub.RunOptions{
		Content: `
hello-provider:
    error: false
hello-provider@run:
    error: true
`,
	})
}
