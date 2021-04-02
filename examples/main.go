package main

import (
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"

	// import all providers
	_ "github.com/erda-project/erda-infra/examples/helloworld"
	_ "github.com/erda-project/erda-infra/examples/protocol/client"
	_ "github.com/erda-project/erda-infra/examples/protocol/pb"
	_ "github.com/erda-project/erda-infra/providers"
)

func main() {
	hub := servicehub.New()
	hub.Run("examples", "examples.yaml", os.Args...)
}
