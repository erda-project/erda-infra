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

func (p *provider) Run(ctx context.Context) error {
	go func() {
		time.Sleep(10 * time.Second)
		err := p.Lock.Close()
		if err != nil {
			fmt.Println("Close err: ", err)
			return
		}
		err = p.Lock.Close()
		if err != nil {
			fmt.Println("Close err: ", err)
			return
		}
	}()
	lock := p.Lock
	sleep := func(d time.Duration) bool {
		select {
		case <-time.After(d):
		case <-ctx.Done():
			return false
		}
		return true
	}
	doTaskInLock := func(prefix string) bool {
		err := lock.Lock(ctx)
		if err != nil {
			fmt.Println("Lock err: ", err)
			return false
		}
		defer func() {
			err := lock.Unlock(context.TODO())
			if err != nil {
				fmt.Println("Unlock err: ", err)
			}
		}()
		fmt.Println(prefix+" {", time.Now())
		if !sleep(1 * time.Second) {
			return false
		}
		fmt.Println(prefix+"    ", time.Now())
		if !sleep(1 * time.Second) {
			return false
		}
		fmt.Println(prefix+" }", time.Now())
		return true
	}
	var wg sync.WaitGroup
	for _, elem := range []string{"A", "B", "C"} {
		wg.Add(1)
		go func(prefix string) {
			defer wg.Done()
			for doTaskInLock(prefix) {
			}
		}(elem)
	}
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
// INFO[2021-09-13 17:13:51.818] provider etcd initialized
// INFO[2021-09-13 17:13:51.818] provider etcd-mutex (depends services: [etcd]) initialized
// INFO[2021-09-13 17:13:51.818] provider example (depends services: [etcd-mutex], providers: [etcd-mutex etcd-mutex]) initialized
// INFO[2021-09-13 17:13:51.818] signals to quit: [hangup interrupt terminated quit]
// INFO[2021-09-13 17:13:51.820] provider example running ...
// A { 2021-09-13 17:13:51.846363 +0800 CST m=+0.082408900
// A     2021-09-13 17:13:52.847656 +0800 CST m=+1.083695286
// A } 2021-09-13 17:13:53.848607 +0800 CST m=+2.084640009
// C { 2021-09-13 17:13:53.872487 +0800 CST m=+2.108520071
// C     2021-09-13 17:13:54.87288 +0800 CST m=+3.108907011
// C } 2021-09-13 17:13:55.874864 +0800 CST m=+4.110884363
// B { 2021-09-13 17:13:55.898284 +0800 CST m=+4.134303984
// B     2021-09-13 17:13:56.902232 +0800 CST m=+5.138246123
// B } 2021-09-13 17:13:57.904445 +0800 CST m=+6.140452108
// A { 2021-09-13 17:13:57.925672 +0800 CST m=+6.161679201
// A     2021-09-13 17:13:58.930195 +0800 CST m=+7.166196424
// A } 2021-09-13 17:13:59.933883 +0800 CST m=+8.169877844
// C { 2021-09-13 17:13:59.956048 +0800 CST m=+8.192042605
// C     2021-09-13 17:14:00.960022 +0800 CST m=+9.196010725
// Lock err:  mutex closed
// Lock err:  mutex closed
// Lock err:  mutex closed
// INFO[2021-09-13 17:14:01.963] provider example Run exit
