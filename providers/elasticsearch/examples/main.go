// Author: recallsong
// Email: songruiguo@qq.com

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
func (d *define) Service() []string { return []string{"example"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{"elasticsearch"} }

// Describe information about this provider
func (d *define) Description() string { return "example" }

// Return an instance representing the configuration
func (d *define) Config() interface{} { return &config{} }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type config struct{}

type provider struct {
	C      *config
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
