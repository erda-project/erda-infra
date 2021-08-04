# servicehub

The *servicehub.Hub* is Service Manager, which manages the startup, initialization, dependency, and shutdown of services.

Provider provide one or more services, and implement the *servicehub.Provider* interface to provide services.

The *servicehub.Hub* manages all providers registered by function *servicehub.RegisterProvider* .

## Example
The configuration file *examples.yaml*
```yaml
hello-provider:
    message: "hello world"
```

The code file *main.go*
```go
package main

import (
	"context"
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

type config struct {
	Message string `file:"message" flag:"msg" default:"hi" desc:"message to show" env:"HELLO_MESSAGE"`
}

type provider struct {
	Cfg *config
	Log logs.Logger
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.Log.Info("message: ", p.Cfg.Message)
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	p.Log.Info("hello provider is running...")
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			p.Log.Info("do something...")
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.Register("hello-provider", &servicehub.Spec{
		Services:    []string{"hello"},
		Description: "hello for example",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
```

Output:
```sh
➜ go run main.go
INFO[2021-03-23 11:55:34.830] message: hello world                          module=hello-provider
INFO[2021-03-23 11:55:34.830] provider hello-provider initialized          
INFO[2021-03-23 11:55:34.830] signals to quit:[hangup interrupt terminated quit] 
INFO[2021-03-23 11:55:34.831] hello provider is running...                  module=hello-provider
INFO[2021-03-23 11:55:37.831] do something...                               module=hello-provider
INFO[2021-03-23 11:55:40.835] do something...                               module=hello-provider
^C
INFO[2021-03-23 11:55:41.624] provider hello-provider Run exit  
```

[Example details](./examples/run/main.go)

[More Examples](./examples/)

## Reading Config
Support the following ways to read config, the priority from low to high is:
* default Tag In Struct
* System Environment Variable
* .env File Environment Variable
* Config File
* Flag

Supports file formats:
* yaml、yml
* json
* hcl
* toml
* ...

