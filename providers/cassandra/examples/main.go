// Copyright (c) 2021 Terminus, Inc.

// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.

// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.

// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/cassandra"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Services() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"cassandra"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	Client cassandra.Interface
}

func (p *provider) Init(ctx servicehub.Context) error {
	session, err := p.Client.Session(&cassandra.SessionConfig{
		Keyspace:    cassandra.KeyspaceConfig{
			Name: "system",
		},
		Consistency: "LOCAL_ONE",
	})
	if err != nil {
		return err
	}
	meta, err := session.KeyspaceMetadata("system")
	if err != nil {
		return err
	}
	fmt.Printf("keyspace name: %s\n", meta.Name)

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
// INFO[2021-04-08 17:49:03.504] provider cassandra initialized
// keyspace name: system
// INFO[2021-04-08 17:49:05.031] provider example (depends [cassandra]) initialized
// INFO[2021-04-08 17:49:05.031] signals to quit: [hangup interrupt terminated quit]