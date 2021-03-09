// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Service() []string { return []string{"hello"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{} }

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
	Message   string    `file:"message" flag:"msg" default:"hi" desc:"message to show"`
	SubConfig subConfig `file:"sub"`
}

type subConfig struct {
	Name string `file:"name" flag:"hello_name" default:"recallsong" desc:"name to show"`
}

type provider struct {
	C       *config
	L       logs.Logger
	closeCh chan struct{}
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.L.Info("message: ", p.C.Message)
	p.closeCh = make(chan struct{})
	return nil
}

func (p *provider) Start() error {
	p.L.Info("now hello provider is running...")
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			p.L.Info("do something...")
		case <-p.closeCh:
			return nil
		}
	}
}

func (p *provider) Close() error {
	p.L.Info("now hello provider is closing...")
	close(p.closeCh)
	return nil
}

func init() {
	servicehub.RegisterProvider("hello", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
