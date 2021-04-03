package main

import (
	"context"
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/examples/protocol/client"
	"github.com/erda-project/erda-infra/examples/protocol/pb"

	// import all providers
	_ "github.com/erda-project/erda-infra/providers"
)

type config struct {
	Name string `file:"name" default:"recallsong"`
}

type provider struct {
	Cfg    *config
	Log    logs.Logger
	Client client.Client // autowired
}

func (p *provider) Run(ctx context.Context) error {
	p.Log.Info("client example is running ...")
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			resp, err := p.Client.GreeterService().SayHello(context.TODO(), &pb.HelloRequest{
				Name: p.Cfg.Name,
			})
			if err != nil {
				p.Log.Error(err)
			}
			p.Log.Info(resp)
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.Register("client-example", &servicehub.Spec{
		Services:     []string{"client-example"},
		Description:  "this is client example",
		Dependencies: []string{"erda.infra.example-client"},
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
	hub.Run("client", "client.yaml", os.Args...)
}
