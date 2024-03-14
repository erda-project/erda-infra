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

package mysqlxorm

import (
	"context"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql client driver package
	"xorm.io/xorm"
	xormlog "xorm.io/xorm/log"
	"xorm.io/xorm/names"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/mysqldriver"
)

// Interface .
type Interface interface {
	DB() *xorm.Engine
	NewSession(ops ...SessionOption) *Session
	Close() error
	GetCloseState() bool
	DataSourceName() string
}

var (
	interfaceType = reflect.TypeOf((*Interface)(nil)).Elem()
	xormType      = reflect.TypeOf((*xorm.Engine)(nil))
)

type config struct {
	MySQLURL            string        `file:"url" env:"MYSQL_URL"`
	MySQLHost           string        `file:"host" env:"MYSQL_HOST" default:"localhost"`
	MySQLPort           string        `file:"port" env:"MYSQL_PORT" default:"3306"`
	MySQLUsername       string        `file:"username" env:"MYSQL_USERNAME" default:"root"`
	MySQLPassword       string        `file:"password" env:"MYSQL_PASSWORD" default:""`
	MySQLDatabase       string        `file:"database" env:"MYSQL_DATABASE"`
	MySQLMaxIdleConns   uint64        `file:"max_idle_conns" env:"MYSQL_MAXIDLECONNS" default:"10"`
	MySQLMaxOpenConns   uint64        `file:"max_open_conns" env:"MYSQL_MAXOPENCONNS" default:"20"`
	MySQLMaxLifeTime    time.Duration `file:"max_lifetime" env:"MYSQL_MAXLIFETIME" default:"30m"`
	MySQLShowSQL        bool          `file:"show_sql" env:"MYSQL_SHOW_SQL" default:"false"`
	MySQLProperties     string        `file:"properties" env:"MYSQL_PROPERTIES" default:"charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"`
	MySQLTLS            string        `file:"tls" env:"MYSQL_TLS"`
	MySQLCaCertPath     string        `file:"ca_cert_path" env:"MYSQL_CACERTPATH"`
	MySQLClientCertPath string        `file:"client_cert_path" env:"MYSQL_CLIENTCERTPATH"`
	MySQLClientKeyPath  string        `file:"client_key_path" env:"MYSQL_CLIENTKEYPATH"`

	MySQLPingWhenInit   bool   `file:"ping_when_init" env:"MYSQL_PING_WHEN_INIT" default:"true"`
	MySQLPingTimeoutSec uint64 `file:"ping_timeout_sec" env:"MYSQL_PING_TIMEOUT_SEC" default:"10"`
}

func (c *config) url() string {
	if c.MySQLURL != "" {
		return c.MySQLURL
	}

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		c.MySQLUsername, c.MySQLPassword, c.MySQLHost, c.MySQLPort, c.MySQLDatabase, c.MySQLProperties)
	if c.MySQLTLS != "" {
		url = fmt.Sprintf("%v&tls=%s", url, c.MySQLTLS)
	}
	return url
}

// provider .
type provider struct {
	Cfg        *config
	Log        logs.Logger
	db         *xorm.Engine
	closeState int32
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	err := mysqldriver.OpenTLS(p.Cfg.MySQLTLS, p.Cfg.MySQLCaCertPath, p.Cfg.MySQLClientCertPath, p.Cfg.MySQLClientKeyPath)
	if err != nil {
		return err
	}

	db, err := xorm.NewEngine("mysql", p.Cfg.url())
	if err != nil {
		return fmt.Errorf("failed to connect to mysql server, err: %v", err)
	}

	db.SetMapper(names.GonicMapper{})
	if p.Cfg.MySQLShowSQL {
		db.ShowSQL(true)
		db.SetLogLevel(xormlog.LOG_DEBUG)
	}

	// connection pool
	db.SetMaxIdleConns(int(p.Cfg.MySQLMaxIdleConns))
	db.SetMaxOpenConns(int(p.Cfg.MySQLMaxOpenConns))
	db.SetConnMaxLifetime(p.Cfg.MySQLMaxLifeTime)
	db.SetDisableGlobalCache(true)

	// ping when init
	if p.Cfg.MySQLPingWhenInit {
		ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.Cfg.MySQLPingTimeoutSec))
		defer cancel()
		if err := db.PingContext(ctxForPing); err != nil {
			return err
		}
	}

	p.db = db
	return nil
}

func (p *provider) DB() *xorm.Engine { return p.db }

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Service() == "mysql-xorm-client" || ctx.Type() == xormType {
		return p.db
	}
	return p
}

func (p *provider) DataSourceName() string {
	return p.DB().DataSourceName()
}

func (p *provider) Close() error {
	if atomic.CompareAndSwapInt32(&p.closeState, 0, 1) {
		return p.db.Close()
	}
	return nil
}

func (p *provider) GetCloseState() bool {
	return atomic.LoadInt32(&p.closeState) == 1
}

func init() {
	servicehub.Register("mysql-xorm", &servicehub.Spec{
		Services: []string{"mysql-xorm", "mysql-xorm-client"},
		Types: []reflect.Type{
			interfaceType, xormType,
		},
		Description: "mysql-xorm",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
