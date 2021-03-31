// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	_ "github.com/erda-project/erda-infra/providers/grpcserver"
)

type provider struct {
}

func (p *provider) Init(ctx servicehub.Context) error {
	return nil
}

func init() {
	servicehub.Register("examples", &servicehub.Spec{
		Services:     []string{"hello"},
		Dependencies: []string{"grpc-server"},
		Description:  "hello for example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
