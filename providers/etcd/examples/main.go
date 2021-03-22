// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"fmt"
	"os"

	"github.com/coreos/etcd/clientv3"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/etcd"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Service() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"etcd"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	ETCD   etcd.Interface   // autowired
	Client *clientv3.Client // autowired
}

func (p *provider) Init(ctx servicehub.Context) error {
	fmt.Println(p.ETCD)
	fmt.Println(p.Client)
	return nil
}

func init() {
	servicehub.RegisterProvider("example", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}

// OUTPUT:
// INFO[2021-03-18 16:26:31.145] provider etcd initialized
// &{0xc00007eaf0 0xc00000ff40 0xc0002029c0 0xc0002cc180}
// &{0xc0002def60 0xc0002def90 0xc0002dac80 0xc00007f400 0xc0002df080 0xc0002df0b0 0xc0002ac700 {[https://127.0.0.1:2379] 0 10000000000 0 0 0 0 0xc0002cc180   false [] <nil> <nil> false} 0xc0002f91e0 0xc0002e6660 0xc000039b60 0xc00007b3c0 0x10e6070   <nil> [{false} {2097152} {2147483647}] 0xc0002e6600}
// INFO[2021-03-18 16:26:31.145] provider example (depends [etcd]) initialized
// INFO[2021-03-18 16:26:31.145] signals to quit:[hangup interrupt terminated quit]
