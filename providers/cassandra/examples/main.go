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
	"github.com/erda-project/erda-infra/providers/cassandra"
)

type config struct {
	Session cassandra.SessionConfig `file:"session"`
}

type provider struct {
	Client cassandra.Interface
	Logger logs.Logger
	Cfg    *config

	s *cassandra.Session
}

func (p *provider) Init(ctx servicehub.Context) error {
	session, err := p.Client.NewSession(&p.Cfg.Session)
	p.s = session
	if err != nil {
		return err
	}

	return nil
}

func (p *provider) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			data, err := p.s.Session().Query("SELECT cql_version FROM system.local").Iter().SliceMap()
			if err != nil {
				p.Logger.Errorf("query error: %s", err)
				continue
			}
			p.Logger.Infof("data: %+v\n", data)
		}
	}
}

func init() {
	servicehub.Register("example", &servicehub.Spec{
		Services:     []string{"example"},
		Dependencies: []string{"cassandra"},
		Description:  "example",
		ConfigFunc: func() interface{} {
			return &config{}
		},
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
// INFO[2021-04-08 17:49:03.504] provider cassandra initialized
// keyspace name: system
// INFO[2021-04-08 17:49:05.031] provider example (depends [cassandra]) initialized
// INFO[2021-04-08 17:49:05.031] signals to quit: [hangup interrupt terminated quit]
