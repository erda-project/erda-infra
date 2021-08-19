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

package election

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/erda-project/erda-infra/base/logs/logrusx"
)

var waitTime = 300 * time.Millisecond

func TestElection(t *testing.T) {
	cfg := integration.ClusterConfig{Size: 1}
	clus := integration.NewClusterV3(nil, &cfg)
	endpoints := []string{clus.Client(0).Endpoints()[0]}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		t.FailNow()
	}

	primary := &provider{Cfg: &config{Prefix: "/election"}, Log: logrusx.New(), Client: cli}
	secondary := &provider{Cfg: &config{Prefix: "/election"}, Log: logrusx.New(), Client: cli}

	leaderSet := int32(0)
	primary.Init(nil)
	primary.OnLeader(func(ctx context.Context) {
		atomic.AddInt32(&leaderSet, 1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	go primary.Run(ctx)
	time.Sleep(waitTime)
	if !primary.IsLeader() || atomic.LoadInt32(&leaderSet) == 0 {
		t.Fatalf("no leader is selected")
	}

	leader, err := primary.Leader()
	if leader == nil || err != nil {
		t.Fatalf("leader function failed")
	}

	secondary.Init(nil)
	secondaryCtx, secondaryCancel := context.WithCancel(context.Background())
	defer secondaryCancel()
	go secondary.Run(secondaryCtx)
	time.Sleep(waitTime)
	if secondary.IsLeader() {
		t.Fatalf("secondary should not be elected")
	}

	leader, err = secondary.Leader()
	if leader == nil || err != nil {
		t.Fatalf("leader function failed")
	}

	// primary exit
	cancel()

	time.Sleep(waitTime)
	if !secondary.IsLeader() {
		t.Fatalf("secondary should be elected")
	}
}

func TestElectionOnClusterReboot(t *testing.T) {
	cfg := integration.ClusterConfig{Size: 1}
	clus := integration.NewClusterV3(nil, &cfg)
	endpoints := []string{clus.Client(0).Endpoints()[0]}
	cli, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
	if err != nil {
		t.FailNow()
	}

	primary := &provider{Cfg: &config{Prefix: "/election", NodeID: "primary"}, Log: logrusx.New(), Client: cli}
	secondary := &provider{Cfg: &config{Prefix: "/election", NodeID: "secondary"}, Log: logrusx.New(), Client: cli}

	primary.Init(nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go primary.Run(ctx)

	time.Sleep(150 * time.Millisecond)
	if !primary.IsLeader() {
		t.Fatalf("no leader is selected")
	}

	secondary.Init(nil)
	secondaryCtx, secondaryCancel := context.WithCancel(context.Background())
	defer secondaryCancel()
	go secondary.Run(secondaryCtx)
	time.Sleep(waitTime)

	// simulate secondary losing connection
	secondary.lock.Lock()
	secondary.session.Close()
	secondary.lock.Unlock()
	// primary quit
	cancel()

	time.Sleep(waitTime)
	if !secondary.IsLeader() {
		t.Fatalf("secondary should be leader now")
	}
}
