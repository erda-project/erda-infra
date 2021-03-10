// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	_ "github.com/erda-project/erda-infra/providers/health"
	_ "github.com/erda-project/erda-infra/providers/httpserver"
)

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
