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

package clickhouse

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	ck "github.com/ClickHouse/clickhouse-go/v2"
	ckdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
)

type Interface interface {
	Client() ckdriver.Conn
}

var (
	interfaceType  = reflect.TypeOf((*Interface)(nil)).Elem()
	nativeConnType = reflect.TypeOf((*ckdriver.Conn)(nil)).Elem()
)

type config struct {
	Addr             string        `file:"addr" default:"localhost:9000"`
	Username         string        `file:"username" default:"default"`
	Password         string        `file:"password"`
	Database         string        `file:"database"`
	MaxIdleConns     int           `file:"max_idle_conns" default:"5"`
	MaxOpenConns     int           `file:"max_open_conns" default:"10"`
	ConnMaxLifeTime  time.Duration `file:"conn_max_lifetime" default:"1h"`
	ConnOpenStrategy uint8         `file:"conn_open_strategy" default:"0"`
	Debug            bool          `file:"debug"`
}

type provider struct {
	Cfg *config
	Log logs.Logger

	nativeConn ckdriver.Conn
}

func (p *provider) Init(ctx servicehub.Context) error {
	options := &ck.Options{
		Addr: strings.Split(p.Cfg.Addr, ","),
		Auth: ck.Auth{
			Database: p.Cfg.Database,
			Username: p.Cfg.Username,
			Password: p.Cfg.Password,
		},
		Debug:            p.Cfg.Debug,
		MaxIdleConns:     p.Cfg.MaxIdleConns,
		MaxOpenConns:     p.Cfg.MaxOpenConns,
		ConnMaxLifetime:  p.Cfg.ConnMaxLifeTime,
		ConnOpenStrategy: ck.ConnOpenStrategy(p.Cfg.ConnOpenStrategy),
	}

	conn, err := ck.Open(options)
	if err != nil {
		return fmt.Errorf("fail to connect clickhouse: %s", err)
	}
	p.nativeConn = conn

	return nil
}

func (p *provider) Client() ckdriver.Conn {
	return p.nativeConn
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Service() == "clickhouse-client" || ctx.Type() == nativeConnType {
		return p.nativeConn
	}
	return p
}

func init() {
	servicehub.Register("clickhouse", &servicehub.Spec{
		Services:    []string{"clickhouse", "clickhouse-client"},
		Types:       []reflect.Type{interfaceType, nativeConnType},
		Description: "clickhouse client",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
