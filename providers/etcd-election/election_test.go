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

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/tests/v3/integration"

	"github.com/erda-project/erda-infra/base/logs/logrusx"
)

// waitTime defines the duration to wait for leader election processes.
var waitTime = 300 * time.Millisecond

// setupCluster initializes a new etcd cluster for testing and returns the cluster and client.
func setupCluster(t *testing.T) (*integration.ClusterV3, *clientv3.Client) {
	integration.BeforeTest(t)
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{
		Size:      3,
		UseBridge: true,
	})
	return cluster, cluster.RandClient()
}

// TestElection tests the basic leader election process.
func TestElection(t *testing.T) {
	cluster, cli := setupCluster(t)
	defer cluster.Terminate(t)

	primary := &provider{Cfg: &config{Prefix: "/election"}, Log: logrusx.New(), Client: cli}
	secondary := &provider{Cfg: &config{Prefix: "/election"}, Log: logrusx.New(), Client: cli}

	leaderSet := int32(0)
	if err := primary.Init(nil); err != nil {
		t.Fatal(err)
	}
	primary.OnLeader(func(ctx context.Context) {
		atomic.AddInt32(&leaderSet, 1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	go primary.Run(ctx)
	time.Sleep(waitTime)

	// Verify primary is elected as leader
	if !primary.IsLeader() || atomic.LoadInt32(&leaderSet) == 0 {
		t.Fatal("no leader is selected")
	}

	leader, err := primary.Leader()
	if leader == nil || err != nil {
		t.Fatalf("leader function failed: %v", err)
	}

	// Initialize secondary and verify it is not elected as leader
	secondary.Init(nil)
	secondaryCtx, secondaryCancel := context.WithCancel(context.Background())
	defer secondaryCancel()
	go secondary.Run(secondaryCtx)
	time.Sleep(waitTime)

	if secondary.IsLeader() {
		t.Fatal("secondary should not be elected")
	}

	leader, err = secondary.Leader()
	if leader == nil || err != nil {
		t.Fatalf("leader function failed: %v", err)
	}

	// Simulate primary exit and verify secondary becomes leader
	cancel()
	time.Sleep(waitTime)
	if !secondary.IsLeader() {
		t.Fatal("secondary should be elected")
	}
}

// TestElectionOnClusterReboot tests leader election after cluster reboot.
func TestElectionOnClusterReboot(t *testing.T) {
	cluster, cli := setupCluster(t)
	defer cluster.Terminate(t)

	primary := &provider{Cfg: &config{Prefix: "/election", NodeID: "primary"}, Log: logrusx.New(), Client: cli}
	secondary := &provider{Cfg: &config{Prefix: "/election", NodeID: "secondary"}, Log: logrusx.New(), Client: cli}

	primary.Init(nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go primary.Run(ctx)
	time.Sleep(waitTime)

	// Verify primary is elected as leader
	if !primary.IsLeader() {
		t.Fatalf("no leader is selected")
	}

	secondary.Init(nil)
	secondaryCtx, secondaryCancel := context.WithCancel(context.Background())
	defer secondaryCancel()
	go secondary.Run(secondaryCtx)
	time.Sleep(waitTime)

	// Simulate secondary losing connection
	secondary.lock.Lock()
	secondary.session.Close()
	secondary.lock.Unlock()

	// Simulate primary exit and verify secondary becomes leader
	cancel()
	time.Sleep(waitTime)
	if !secondary.IsLeader() {
		t.Fatalf("secondary should be leader now")
	}
}
