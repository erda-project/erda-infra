// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"os"

	_ "github.com/erda-project/erda-infra/providers/httpserver"
	_ "github.com/erda-project/erda-infra/providers/pprof"
	"github.com/erda-project/erda-infra/base/servicehub"
)

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
