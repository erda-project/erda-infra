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
var waitTime = time.Second

// setupCluster initializes a new etcd cluster for testing and returns the cluster and client.
func setupCluster(t *testing.T) (*integration.ClusterV3, *clientv3.Client) {
	integration.BeforeTest(t)
	cluster := integration.NewClusterV3(t, &integration.ClusterConfig{
		Size:      3,
		UseBridge: true,
	})
	return cluster, cluster.RandClient()
}

// initProvider initializes a provider with the given client and prefix.
func initProvider(cli *clientv3.Client, prefix string) *provider {
	return &provider{Cfg: &config{Prefix: prefix}, Log: logrusx.New(), Client: cli}
}

// runProvider runs the provider in a separate goroutine.
func runProvider(p *provider, ctx context.Context) {
	go p.Run(ctx)
	time.Sleep(waitTime)
}

// verifyLeader verifies if the provider is the leader.
func verifyLeader(t *testing.T, p *provider, expected bool) {
	if p.IsLeader() != expected {
		t.Fatalf("expected leader: %v, but got: %v", expected, p.IsLeader())
	}
}

// TestElection tests the basic leader election process.
func TestElection(t *testing.T) {
	cluster, cli := setupCluster(t)
	defer cluster.Terminate(t)

	tests := []struct {
		name        string
		path        string
		primaryInit func(*provider)
		runTest     func(*testing.T, *provider, *provider)
	}{
		{
			name: "basic leader election",
			path: "/election/case1",
			primaryInit: func(primary *provider) {
				leaderSet := int32(0)
				primary.OnLeader(func(ctx context.Context) {
					atomic.AddInt32(&leaderSet, 1)
				})
			},
			runTest: func(t *testing.T, primary, secondary *provider) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				runProvider(primary, ctx)
				verifyLeader(t, primary, true)

				secondaryCtx, secondaryCancel := context.WithCancel(context.Background())
				defer secondaryCancel()
				runProvider(secondary, secondaryCtx)
				verifyLeader(t, secondary, false)

				cancel()
				time.Sleep(waitTime)
				verifyLeader(t, secondary, true)
			},
		},
		{
			name:        "leader election on cluster reboot",
			path:        "/election/case2",
			primaryInit: func(primary *provider) {},
			runTest: func(t *testing.T, primary, secondary *provider) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				runProvider(primary, ctx)
				verifyLeader(t, primary, true)

				secondaryCtx, secondaryCancel := context.WithCancel(context.Background())
				defer secondaryCancel()
				runProvider(secondary, secondaryCtx)

				secondary.lock.Lock()
				if err := secondary.session.Close(); err != nil {
					t.Fatalf("failed to close session, %v", err)
				}
				secondary.lock.Unlock()

				cancel()
				time.Sleep(waitTime)
				verifyLeader(t, secondary, true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primary := initProvider(cli, tt.path)
			secondary := initProvider(cli, tt.path)

			if err := primary.Init(nil); err != nil {
				t.Fatal(err)
			}
			tt.primaryInit(primary)
			if err := secondary.Init(nil); err != nil {
				t.Fatal(err)
			}

			tt.runTest(t, primary, secondary)
		})
	}
}
