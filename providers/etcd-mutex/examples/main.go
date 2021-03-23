// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/erda-project/erda-infra/base/servicehub"
	mutex "github.com/erda-project/erda-infra/providers/etcd-mutex"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Services() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"etcd-mutex"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	Mutex mutex.Interface // autowired
	Lock  mutex.Mutex     `mutex-key:"test-key"` // autowired
}

func (p *provider) Init(ctx servicehub.Context) error {
	fmt.Println(p.Mutex)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer p.Lock.Unlock(context.TODO())
		p.Lock.Lock(context.TODO())
		fmt.Println("A", time.Now())
		time.Sleep(5 * time.Second)
	}()
	go func() {
		defer wg.Done()
		defer p.Lock.Unlock(context.TODO())
		p.Lock.Lock(context.TODO())
		fmt.Println("B", time.Now())
		time.Sleep(5 * time.Second)
	}()
	wg.Wait()
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
// INFO[2021-03-18 17:18:46.578] provider etcd initialized
// INFO[2021-03-18 17:18:46.578] provider etcd-mutex (depends [etcd]) initialized
// &{0xc00000ffe0 0xc0002f9800 0xc0002b8630 map[] 0xc0002f6070}
// B 2021-03-18 17:18:46.772563 +0800 CST m=+0.202808470
// A 2021-03-18 17:18:51.804971 +0800 CST m=+5.235190119
// INFO[2021-03-18 17:18:56.852] provider example (depends [etcd-mutex]) initialized
// INFO[2021-03-18 17:18:56.852] signals to quit:[hangup interrupt terminated quit]
