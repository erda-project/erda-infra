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

package modcom

import (
	"os"
	"path"
	"strings"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/base/version"
	"github.com/erda-project/erda-infra/modcom/api"
	"github.com/erda-project/erda-infra/providers/i18n"
	"github.com/recallsong/go-utils/config"
	uuid "github.com/satori/go.uuid"
)

var instanceID = uuid.NewV4().String()

// InstanceID .
func InstanceID() string { return instanceID }

// Env .
func Env() {
	config.LoadEnvFile()
}

// GetEnv get environment with default value
func GetEnv(key, def string) string {
	v := os.Getenv(key)
	if len(v) > 0 {
		return v
	}
	return def
}

func loadModuleEnvFile(dir string) {
	path := path.Join(dir, ".env")
	config.LoadEnvFileWithPath(path, false)
}

func prepare() {
	version.PrintIfCommand()
	Env()
	for _, fn := range initializers {
		fn()
	}
}

var initializers []func()

// RegisterInitializer .
func RegisterInitializer(fn func()) {
	initializers = append(initializers, fn)
}

// Hub global variable
var Hub *servicehub.Hub

var listener = &servicehub.DefaultListener{
	BeforeInitFunc: func(h *servicehub.Hub, config map[string]interface{}) error {
		if _, ok := config["i18n"]; !ok {
			config["i18n"] = nil // i18n is required
		}
		return nil
	},
	AfterInitFunc: func(h *servicehub.Hub) error {
		api.I18n = h.Service("i18n").(i18n.I18n)
		return nil
	},
}

// Run .
func Run(cfg string) {
	prepare()
	Hub := servicehub.New(servicehub.WithListener(listener))
	Hub.Run("", cfg, os.Args...)
}

// RunWithCfgDir .
func RunWithCfgDir(dir, name string) {
	prepare()
	name = GetEnv("CONFIG_NAME", name)
	dir = strings.TrimRight(dir, "/")
	os.Setenv("CONFIG_PATH", dir)
	loadModuleEnvFile(dir)
	cfg := path.Join(dir, name+GetEnv("CONFIG_SUFFIX", ".yaml"))

	// create and run service hub
	Hub := servicehub.New(servicehub.WithListener(listener))
	Hub.Run("", cfg, os.Args...)
}
