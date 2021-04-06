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
	"time"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/elasticsearch"
	"github.com/olivere/elastic"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Services() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"elasticsearch"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	ES     elasticsearch.Interface // autowired
	Client *elastic.Client         // autowired
}

func (p *provider) Init(ctx servicehub.Context) error {
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := p.Client.CatIndices().Do(context)
	if err != nil {
		return err
	}
	for _, item := range resp {
		fmt.Println(item.Index)
	}
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
// NFO[2021-03-18 16:14:17.725] provider elasticsearch initialized
// spot-elasticsearch_http-full_cluster-1615939200000
// spot-elasticsearch_transport-full_cluster-1615766400000
// INFO[2021-03-18 16:14:19.802] provider example (depends [elasticsearch]) initialized
// INFO[2021-03-18 16:14:19.802] signals to quit:[hangup interrupt terminated quit]
