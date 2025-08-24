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
	"path/filepath"
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
	fmt.Println("running ...")
	go func() {
		time.Sleep(10000 * time.Second)
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
	//ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	mu, err := p.Mutex.New(ctx, "keyAAAA")
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := mu.Lock(ctx)
		if err != nil {
			return
		}
		defer mu.Unlock(ctx)
		for i := 0; i < 10; i++ {
			fmt.Printf("A wait %vs\n", i+1)
			time.Sleep(1 * time.Second)
		}
		fmt.Println("AAA===>")
	}()
	time.Sleep(1 * time.Second)
	ctx, cancelFunc = context.WithCancel(context.Background())
	defer cancelFunc()
	mu, err = p.Mutex.New(ctx, "keyAAAA")
	if err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := mu.Lock(ctx)
		if err != nil {
			return
		}
		defer mu.Unlock(ctx)
		time.Sleep(5 * time.Second)
		fmt.Println("BBB===>")
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
	dir, _ := os.Getwd()
	hub := servicehub.New()
	hub.Run("examples", filepath.Join(dir, "providers", "etcd-mutex", "examples", "examples.yaml"), os.Args...)
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
