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
	servicehub.Register("example", &servicehub.Spec{
		Services:     []string{"example"},
		Dependencies: []string{"etcd-mutex"},
		Description:  "example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
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
