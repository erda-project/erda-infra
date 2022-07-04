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

package v2

import (
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql" // mysql client driver package
	"gorm.io/gorm"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/mysqldriver"
)

var (
	interfaceType = reflect.TypeOf((*Interface)(nil)).Elem()
	gormType      = reflect.TypeOf((*gorm.DB)(nil))
	txType        = reflect.TypeOf((*TX)(nil))
	name          = "gorm.v2"
	spec          = servicehub.Spec{
		Services: []string{"mysql-gorm.v2", "mysql-gorm.v2-client"},
		Types: []reflect.Type{
			interfaceType, gormType,
		},
		Description: "mysql-gorm.v2",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	}
)

func init() {
	servicehub.Register(name, &spec)
}

// Interface .
type Interface interface {
	DB() *gorm.DB
	Q() *TX
	Begin() *TX
}

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
	MySQLDebug          bool          `file:"debug" env:"MYSQL_DEBUG" default:"false"`
	MySQLCharset        string        `file:"charset" env:"MYSQL_CHARSET" default:"utf8mb4"`
	MySQLTLS            string        `file:"tls" env:"MYSQL_TLS"`
	MySQLCaCertPath     string        `file:"ca_cert_path" env:"MYSQL_CACERTPATH"`
	MySQLClientCertPath string        `file:"client_cert_path" env:"MYSQL_CLIENTCERTPATH"`
	MySQLClientKeyPath  string        `file:"client_key_path" env:"MYSQL_CLIENTKEYPATH"`
}

func (c *config) url() string {
	if c.MySQLURL != "" {
		return c.MySQLURL
	}

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.MySQLUsername, c.MySQLPassword, c.MySQLHost, c.MySQLPort, c.MySQLDatabase, c.MySQLCharset)
	if c.MySQLTLS != "" {
		url = fmt.Sprintf("%v&tls=%s", url, c.MySQLTLS)
	}
	return url
}

// provider .
type provider struct {
	Cfg *config
	db  *gorm.DB
	tx  *TX
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	err := mysqldriver.OpenTLS(p.Cfg.MySQLTLS, p.Cfg.MySQLCaCertPath, p.Cfg.MySQLClientCertPath, p.Cfg.MySQLClientKeyPath)
	if err != nil {
		return err
	}

	logrus.WithField("provider", name).Infoln("init")
	db, err := gorm.Open(mysql.Open(p.Cfg.url()))
	if err != nil {
		return fmt.Errorf("fail to connect mysql: %s", err)
	}

	s, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get a conn pool")
	}

	// connection pool
	s.SetMaxIdleConns(int(p.Cfg.MySQLMaxIdleConns))
	s.SetMaxOpenConns(int(p.Cfg.MySQLMaxOpenConns))
	s.SetConnMaxLifetime(p.Cfg.MySQLMaxLifeTime)
	p.db = db
	if p.Cfg.MySQLDebug {
		p.db = p.db.Debug()
	}

	return nil
}

// DB .
func (p *provider) DB() *gorm.DB { return p.db }

func (p *provider) Q() *TX {
	return &TX{
		db:    p.db,
		valid: true,
	}
}

func (p *provider) Begin() *TX {
	return &TX{
		db:    p.db,
		valid: true,
		inTx:  true,
	}
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Service() == "mysql-gorm.v2-client" || ctx.Type() == gormType {
		return p.db
	}
	return p
}
