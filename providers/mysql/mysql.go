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

package mysql

import (
	"fmt"
	"reflect"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	_ "github.com/go-sql-driver/mysql" // mysql client driver package
	"github.com/jinzhu/gorm"
)

// Interface .
type Interface interface {
	DB() *gorm.DB
}

var (
	interfaceType = reflect.TypeOf((*Interface)(nil)).Elem()
	gormType      = reflect.TypeOf((*gorm.DB)(nil))
)

type config struct {
	MySQLURL          string        `file:"url" env:"MYSQL_URL"`
	MySQLHost         string        `file:"host" env:"MYSQL_HOST" default:"localhost"`
	MySQLPort         string        `file:"port" env:"MYSQL_PORT" default:"3306"`
	MySQLUsername     string        `file:"username" env:"MYSQL_USERNAME" default:"root"`
	MySQLPassword     string        `file:"password" env:"MYSQL_PASSWORD" default:""`
	MySQLDatabase     string        `file:"database" env:"MYSQL_DATABASE"`
	MySQLMaxIdleConns uint64        `file:"max_idle_conns" env:"MYSQL_MAXIDLECONNS" default:"1"`
	MySQLMaxOpenConns uint64        `file:"max_open_conns" env:"MYSQL_MAXOPENCONNS" default:"2"`
	MySQLMaxLifeTime  time.Duration `file:"max_lifetime" env:"MYSQL_MAXLIFETIME" default:"30m"`
}

func (c *config) url() string {
	if c.MySQLURL != "" {
		return c.MySQLURL
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.MySQLUsername, c.MySQLPassword, c.MySQLHost, c.MySQLPort, c.MySQLDatabase)
}

type define struct{}

func (d *define) Services() []string { return []string{"mysql", "mysql-client"} }
func (d *define) Types() []reflect.Type {
	return []reflect.Type{
		interfaceType, gormType,
	}
}
func (d *define) Summary() string     { return "mysql" }
func (d *define) Description() string { return d.Summary() }
func (d *define) Config() interface{} { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

// provider .
type provider struct {
	Cfg *config
	Log logs.Logger
	db  *gorm.DB
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	db, err := gorm.Open("mysql", p.Cfg.url())
	if err != nil {
		return fmt.Errorf("fail to connect mysql: %s", err)
	}

	// connection pool
	db.DB().SetMaxIdleConns(int(p.Cfg.MySQLMaxIdleConns))
	db.DB().SetMaxOpenConns(int(p.Cfg.MySQLMaxOpenConns))
	db.DB().SetConnMaxLifetime(p.Cfg.MySQLMaxLifeTime)
	p.db = db
	return nil
}

func (p *provider) DB() *gorm.DB { return p.db }

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Service() == "mysql-client" || ctx.Type() == gormType {
		return p.db
	}
	return p
}

func init() {
	servicehub.RegisterProvider("mysql", &define{})
}
