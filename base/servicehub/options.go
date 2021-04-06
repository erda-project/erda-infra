// Copyright 2021 Terminus
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

package servicehub

import "github.com/erda-project/erda-infra/base/logs"

// Option .
type Option func(hub *Hub)

func processOptions(hub *Hub, opt interface{}) {
	if fn, ok := opt.(Option); ok {
		fn(hub)
	}
}

// WithLogger .
func WithLogger(logger logs.Logger) interface{} {
	return Option(func(hub *Hub) {
		hub.logger = logger
	})
}

// Listener .
type Listener interface {
	BeforeInitialization(h *Hub, config map[string]interface{}) error
	AfterInitialization(h *Hub) error
	BeforeExit(h *Hub, err error) error
}

// WithListener .
func WithListener(l Listener) interface{} {
	return Option(func(hub *Hub) {
		hub.listeners = append(hub.listeners, l)
	})
}

// DefaultListener .
type DefaultListener struct {
	BeforeInitFunc func(h *Hub, config map[string]interface{}) error
	AfterInitFunc  func(h *Hub) error
	BeforeExitFunc func(h *Hub, err error) error
}

// BeforeInitialization .
func (l *DefaultListener) BeforeInitialization(h *Hub, config map[string]interface{}) error {
	if l.BeforeInitFunc == nil {
		return nil
	}
	return l.BeforeInitFunc(h, config)
}

// AfterInitialization .
func (l *DefaultListener) AfterInitialization(h *Hub) error {
	if l.AfterInitFunc == nil {
		return nil
	}
	return l.AfterInitFunc(h)
}

// BeforeExit .
func (l *DefaultListener) BeforeExit(h *Hub, err error) error {
	if l.BeforeExitFunc == nil {
		return err
	}
	return l.BeforeExitFunc(h, err)
}
