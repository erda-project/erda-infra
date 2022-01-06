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

package cpregister

import (
	"embed"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
)

// AllExplicitProviderCreatorMap contains all user specified provider.
var AllExplicitProviderCreatorMap = map[string]servicehub.Provider{}

// HubListener .
type HubListener struct {
	ScenarioFSs []embed.FS
}

// NewHubListener .
func NewHubListener(fs ...embed.FS) *HubListener {
	return &HubListener{ScenarioFSs: fs}
}

// BeforeInitialization .
func (l *HubListener) BeforeInitialization(h *servicehub.Hub, config map[string]interface{}) error {
	// auto register explicit component provider firstly
	logrus.Info("auto register component provider to hub.config")
	for providerName, creator := range AllExplicitProviderCreatorMap {
		config[providerName] = creator
		logrus.Infof("auto register component provider to hub.config: %s", providerName)
	}

	// register default protocols from FS
	for _, fs := range l.ScenarioFSs {
		protocol.MustRegisterProtocolsFromFS(fs)
	}

	return nil
}

// AfterInitialization .
func (l *HubListener) AfterInitialization(h *servicehub.Hub) error { return nil }

// AfterStart .
func (l *HubListener) AfterStart(h *servicehub.Hub) error { return nil }

// BeforeExit .
func (l *HubListener) BeforeExit(h *servicehub.Hub, err error) error { return nil }

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05.000"})
}
