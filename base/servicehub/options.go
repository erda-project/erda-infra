// Author: recallsong
// Email: songruiguo@qq.com

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
}

// WithListener .
func WithListener(l Listener) interface{} {
	return Option(func(hub *Hub) {
		hub.listeners = append(hub.listeners, l)
	})
}
