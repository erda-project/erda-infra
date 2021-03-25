// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/kafka"
)

type define struct{}

func (d *define) Service() []string      { return []string{"hello"} }
func (d *define) Dependencies() []string { return []string{"kafka"} }
func (d *define) Config() interface{}    { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type config struct {
	Input kafka.ConsumerConfig `file:"input"`
}

type provider struct {
	Cfg   *config
	Log   logs.Logger
	Kafka kafka.Interface
}

func (p *provider) Run(ctx context.Context) error {
	p.Kafka.NewConsumer(&p.Cfg.Input, p.invoke)
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *provider) invoke(key []byte, value []byte, topic *string, timestamp time.Time) error {
	fmt.Println(string(value))
	return nil
}

func init() {
	servicehub.RegisterProvider("examples", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
