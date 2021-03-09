# servicehub

The *servicehub.Hub* is Service Manager, which manages the startup, initialization, dependency, and shutdown of services.

Provider provide one or more services, and implement the *servicehub.Provider* interface to provide services.

The *servicehub.Hub* manages all providers registered by function *servicehub.RegisterProvider* .

## Example
The configuration file *examples.yaml*
```yaml
hello:
    message: "hello world"
    sub:
        name: "recallsong"
```

The code file *main.go*
```go
package main

import (
	"os"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

// define Represents the definition of provider and provides some information
type define struct{}

// Declare what services the provider provides
func (d *define) Service() []string { return []string{"hello"} }

// Declare which services the provider depends on
func (d *define) Dependencies() []string { return []string{} }

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
	Message   string    `file:"message" flag:"msg" default:"hi" desc:"message to show"`
	SubConfig subConfig `file:"sub"`
}

type subConfig struct {
	Name string `file:"name" flag:"hello_name" default:"recallsong" desc:"name to show"`
}

type provider struct {
	C       *config
	L       logs.Logger
	closeCh chan struct{}
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.L.Info("message: ", p.C.Message)
	p.closeCh = make(chan struct{})
	return nil
}

func (p *provider) Start() error {
	p.L.Info("now hello provider is running...")
	tick := time.Tick(10 * time.Second)
	for {
		select {
		case <-tick:
			p.L.Info("do something...")
		case <-p.closeCh:
			return nil
		}
	}
}

func (p *provider) Close() error {
	p.L.Info("now hello provider is closing...")
	close(p.closeCh)
	return nil
}

func init() {
	servicehub.RegisterProvider("hello", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
```
[Example details](./examples/main.go)

Output:
```sh
➜  examples git:(master) ✗ go run main.go
INFO[2021-03-08 19:04:09.493] message: hello world                          module=hello
INFO[2021-03-08 19:04:09.493] provider hello initialized                   
INFO[2021-03-08 19:04:09.493] signals to quit:[hangup interrupt terminated quit] 
INFO[2021-03-08 19:04:09.493] now hello provider is running...              module=hello
INFO[2021-03-08 19:04:19.496] do something...                               module=hello
INFO[2021-03-08 19:04:29.497] do something...                               module=hello
^C
INFO[2021-03-08 19:04:32.984] now hello provider is closing...              module=hello
INFO[2021-03-08 19:04:32.984] provider hello exit     
```

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
