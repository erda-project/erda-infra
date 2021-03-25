// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	_ "github.com/erda-project/erda-infra/providers/httpserver"
	_ "github.com/erda-project/erda-infra/providers/pprof"
)

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
