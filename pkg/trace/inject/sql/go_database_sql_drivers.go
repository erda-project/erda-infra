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
	"database/sql/driver"
	"log"
	"sync"
	_ "unsafe"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	//go:linkname drivers database/sql.driversMu
	driversMu sync.RWMutex

	//go:linkname drivers database/sql.drivers
	drivers map[string]driver.Driver
)

// WrapDrivers .
func WrapDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	for name, driver := range drivers {
		if _, ok := driver.(*wrappedDriver); ok {
			continue
		}
		if _, ok := driver.(*wrappedDriverContext); ok {
			continue
		}
		drivers[name] = WrapDriver(name, driver)
	}
}

// WrapDriver .
func WrapDriver(name string, d driver.Driver) driver.Driver {
	log.Printf("hook %q database driver", name)
	return wrapDriver(otelsql.WrapDriver(d, name))
}

func init() {
	WrapDrivers()
}
