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

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Service() []string { return []string{"hello"} }

// Describe information about this provider
func (d *define) Description() string { return "hello for example" }

// Return an instance representing the configuration
func (d *define) Config() interface{} { return &config{} }

// Return a provider creator
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type config struct {
	Message string `file:"message" flag:"msg" default:"hi" desc:"message to show" env:"HELLO_MESSAGE"`
}

type provider struct {
	C *config
	L logs.Logger
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.L.Info("message: ", p.C.Message)
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	p.L.Info("hello provider is running...")
	tick := time.NewTicker(3 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			p.L.Info("do something...")
		case <-ctx.Done():
			return nil
		}
	}
}

func init() {
	servicehub.RegisterProvider("hello-provider", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
```

Output:
```sh
➜ go run main.go
INFO[2021-03-18 14:18:26.313] message: hello world                          module=hello-provider
INFO[2021-03-18 14:18:26.313] provider hello-provider initialized          
INFO[2021-03-18 14:18:26.313] signals to quit:[hangup interrupt terminated quit] 
INFO[2021-03-18 14:18:26.314] hello provider is running...                  module=hello-provider
INFO[2021-03-18 14:18:29.315] do something...                               module=hello-provider
INFO[2021-03-18 14:18:32.317] do something...                               module=hello-provider
^C
INFO[2021-03-18 14:18:34.468] provider hello-provider exit   
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

## TODO List
* CLI tools to quick start
* Test Case
