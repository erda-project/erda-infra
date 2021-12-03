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

package sql

import (
	"database/sql"
	"database/sql/driver"
	"sync"

	"github.com/XSAM/otelsql"
	"github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"github.com/erda-project/erda-infra/pkg/trace/inject/hook"
)

//go:noinline
func originalOpen(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

var (
	driversMu sync.Mutex
	drivers   = map[string]string{}
)

//go:noinline
func tracedOpen(driverName, dataSourceName string) (*sql.DB, error) {
	driversMu.Lock()
	if dname, ok := drivers[driverName]; !ok {
		// retrieve the driver implementation we need to wrap with instrumentation
		db, err := originalOpen(driverName, "")
		if err != nil {
			driversMu.Unlock()
			return nil, err
		}
		d := db.Driver()
		if err = db.Close(); err != nil {
			driversMu.Unlock()
			return nil, err
		}
		dname = "otelsql-" + driverName
		sql.Register(dname, wrapDriver(otelsql.WrapDriver(d, driverName, otelsql.WithAttributes(getDriverAttributes(driverName, dataSourceName, d)...))))
		drivers[driverName] = dname
		driverName = dname
	} else {
		driverName = dname
	}
	driversMu.Unlock()
	return originalOpen(driverName, dataSourceName)
}

func getDriverAttributes(driverName, dataSourceName string, d driver.Driver) (attrs []attribute.KeyValue) {
	switch driverName {
	case "mysql":
		cfg, err := mysql.ParseDSN(dataSourceName)
		if err == nil {
			attrs = append(attrs,
				semconv.DBUserKey.String(cfg.User),
				semconv.DBNameKey.String(cfg.DBName),
				attribute.Key("db.host").String(cfg.Addr),
			)
		}
	}
	return attrs
}

func init() {
	hook.Hook(sql.Open, tracedOpen, originalOpen)
}
