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
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	_ "github.com/erda-project/erda-infra/providers/etcd"
	election "github.com/erda-project/erda-infra/providers/etcd-election"
)

type provider struct {
	Log      logs.Logger
	Election election.Interface `autowired:"etcd-election"`
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.Election.OnLeader(p.leaderTask)
	return nil
}

func (p *provider) leaderTask(ctx context.Context) {
	defer p.Log.Info("leader task exit")
	watch := p.Election.Watch(ctx)
	for {
		select {
		case event := <-watch:
			p.Log.Infof("nodes changed: %v", event)
		case <-time.After(3 * time.Second):
			p.Log.Info("leader task doing")
		case <-ctx.Done():
			return
		}
	}
}

func (p *provider) Run(ctx context.Context) error {
	select {
	case <-time.After(10 * time.Second):
		p.Log.Info("resign leader")
		p.Election.ResignLeader()
	case <-ctx.Done():
	}
	return nil
}

func init() {
	servicehub.Register("example", &servicehub.Spec{
		Services:     []string{"example"},
		Dependencies: []string{"etcd"},
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
// provider etcd config:
// INFO[2021-06-23 17:13:26.894] provider etcd initialized
// node id:  eff22718-a19c-4489-aa1c-b874f79c8b60
// INFO[2021-06-23 17:13:26.894] provider etcd-election (depends services: [etcd etcd-client]) initialized
// INFO[2021-06-23 17:13:26.894] provider example (depends services: [etcd etcd-election]) initialized
// INFO[2021-06-23 17:13:26.895] signals to quit: [hangup interrupt terminated quit]
// INFO[2021-06-23 17:13:26.900] provider example running ...
// INFO[2021-06-23 17:13:26.902] provider etcd-election running ...
// INFO[2021-06-23 17:13:29.124] I am leader ! node is "eff22718-a19c-4489-aa1c-b874f79c8b60"  module=etcd-election
// INFO[2021-06-23 17:13:32.129] leader task doing                             module=example
// INFO[2021-06-23 17:13:32.802] nodes changed: {delete {74ee6ac7-1703-43c6-bdc8-df90f9984d5b}}  module=example
// INFO[2021-06-23 17:13:35.807] leader task doing                             module=example
// INFO[2021-06-23 17:13:36.900] resign leader                                 module=example
// INFO[2021-06-23 17:13:36.901] leader task exit                              module=example
// INFO[2021-06-23 17:13:36.968] provider example exit
// INFO[2021-06-23 17:13:42.798] I am leader ! node is "eff22718-a19c-4489-aa1c-b874f79c8b60"  module=etcd-election
// INFO[2021-06-23 17:13:45.800] leader task doing                             module=example
// INFO[2021-06-23 17:13:48.802] leader task doing                             module=example
// INFO[2021-06-23 17:13:51.806] leader task doing                             module=example
// INFO[2021-06-23 17:13:54.497] nodes changed: {delete {}}                    module=example
// INFO[2021-06-23 17:13:57.502] leader task doing                             module=example
// INFO[2021-06-23 17:14:00.506] leader task doing                             module=example
// INFO[2021-06-23 17:14:03.510] leader task doing                             module=example
