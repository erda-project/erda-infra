// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	election "github.com/erda-project/erda-infra/providers/zk-master-election"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Services() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"zk-master-election"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	Election election.Interface // autowired
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.Election.Watch(p.masterChanged)
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
}

func (p *provider) masterChanged(event election.Event) {
	fmt.Println("is master: ", event.IsMaster(), ", is connected: ", event.IsConnected())
}

func init() {
	servicehub.RegisterProvider("example", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}

// OUTPUT:
// INFO[2021-03-18 15:31:36.589] provider zookeeper initialized
// INFO[2021-03-18 15:31:36.589] provider zk-master-election (depends [zookeeper]) initialized
// INFO[2021-03-18 15:31:36.589] provider example (depends [zk-master-election]) initialized
// INFO[2021-03-18 15:31:36.589] signals to quit:[hangup interrupt terminated quit]
// 2021/03/18 15:31:36 connected to 127.0.0.1:2181
// INFO[2021-03-18 15:31:36.682] connected to zookeeper successfully           module=zk-master-election
// 2021/03/18 15:31:36 authenticated: id=105855796925956124, timeout=12000
// 2021/03/18 15:31:36 re-submitting `0` credentials after reconnect
// INFO[2021-03-18 15:31:36.871] election finish, i am slave                   module=zk-master-election
// is master:  true , is connected:  true
// INFO[2021-03-18 15:31:36.929] start watch path "/example/master-node-key"   module=zk-master-election
// ^C
// INFO[2021-03-18 15:31:44.602] provider example exit
// INFO[2021-03-18 15:31:44.602] exit waith path "/example/master-node-key"    module=zk-master-election
// 2021/03/18 15:31:44 recv loop terminated: EOF
// 2021/03/18 15:31:44 send loop terminated: <nil>
// INFO[2021-03-18 15:31:44.634] disconnected zookeeper                        module=zk-master-election
// INFO[2021-03-18 15:31:44.634] provider zk-master-election closed
