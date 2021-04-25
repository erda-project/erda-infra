// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zookeeper

import (
	"reflect"
	"strings"
	"time"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/go-zookeeper/zk"
)

// Interface .
type Interface interface {
	Connect(options ...func(*zk.Conn)) (*zk.Conn, <-chan zk.Event, error)
	SessionTimeout() time.Duration
}

type config struct {
	Addrs          string        `file:"addrs" env:"ZOOKEEPER_ADDR"`
	SessionTimeout time.Duration `file:"session_timeout" default:"3s"`
}

type define struct{}

func (d *define) Services() []string { return []string{"zookeeper"} }
func (d *define) Types() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Interface)(nil)).Elem()}
}
func (d *define) Description() string { return "zookeeper" }
func (d *define) Config() interface{} { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	Cfg *config
}

func (p *provider) Connect(options ...func(*zk.Conn)) (*zk.Conn, <-chan zk.Event, error) {
	return zk.Connect(strings.Split(p.Cfg.Addrs, ","), p.Cfg.SessionTimeout)
}

func (p *provider) SessionTimeout() time.Duration { return p.Cfg.SessionTimeout }

func init() {
	servicehub.RegisterProvider("zookeeper", &define{})
}
