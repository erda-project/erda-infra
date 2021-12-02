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
	"context"
	"database/sql/driver"
)

func wrapDriver(d driver.Driver) driver.Driver {
	drv := &wrappedDriver{
		driver: d,
	}
	if ctx, ok := d.(driver.DriverContext); ok {
		return &wrappedDriverContext{
			wrappedDriver: drv,
			context:       ctx,
		}
	}
	return drv
}

type wrappedDriver struct {
	driver driver.Driver
}

func (d *wrappedDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &wrappedConn{
		Conn: conn,
	}, nil
}

type wrappedDriverContext struct {
	*wrappedDriver
	context driver.DriverContext
}

func (d *wrappedDriverContext) OpenConnector(name string) (driver.Connector, error) {
	rawConnector, err := d.context.OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return &wrappedConnector{
		driver:    d,
		connector: rawConnector,
	}, err
}

type wrappedConnector struct {
	driver    driver.Driver
	connector driver.Connector
}

func (c *wrappedConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &wrappedConn{
		Conn: conn,
	}, nil
}

func (c *wrappedConnector) Driver() driver.Driver {
	return c.driver
}
