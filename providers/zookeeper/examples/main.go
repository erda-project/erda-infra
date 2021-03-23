// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/zookeeper"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Services() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"zookeeper"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	ZooK zookeeper.Interface // autowired
}

func (p *provider) Run(ctx context.Context) error {
	conn, ch, err := p.ZooK.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	for {
		select {
		case event := <-ch:
			// do something
			fmt.Println(event)
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.RegisterProvider("example", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}

// OUTPUT:
// INFO[2021-03-18 15:33:03.721] provider zookeeper initialized
// INFO[2021-03-18 15:33:03.721] provider example (depends [zookeeper]) initialized
// INFO[2021-03-18 15:33:03.721] signals to quit:[hangup interrupt terminated quit]
// {EventSession StateConnecting  <nil> 127.0.0.1:2181}
// 2021/03/18 15:33:04 connected to 127.0.0.1:2181
// {EventSession StateConnected  <nil> 127.0.0.1:2181}
// 2021/03/18 15:33:04 authenticated: id=105855796925956125, timeout=12000
// {EventSession StateHasSession  <nil> 127.0.0.1:2181}
// 2021/03/18 15:33:04 re-submitting `0` credentials after reconnect
// ^C
// 2021/03/18 15:33:09 recv loop terminated: EOF
// 2021/03/18 15:33:09 send loop terminated: <nil>
// INFO[2021-03-18 15:33:09.306] provider example exit
