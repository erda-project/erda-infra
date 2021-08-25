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

package cassandra

import (
	"sync"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/gocql/gocql"
)

type ReconnectionConfig struct {
	Enable        bool          `file:"enable" default:"true"`
	CheckInterval time.Duration `file:"check_interval" default:"10m"`
	CheckTimeout  time.Duration `file:"check_timeout" default:"60s"`
}

type Session struct {
	session *gocql.Session
	mu      sync.RWMutex
	log     logs.Logger
	done    chan struct{}
}

func (s *Session) Session() *gocql.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.session
}

func (s *Session) Close() {
	close(s.done)
}

func (s *Session) updateSession(gs *gocql.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.session = gs
}

func (s *Session) checkAndReconnect(p *provider, cfg *SessionConfig) {
	log := s.log.Sub("reconnection")
	log.Infof("start to check connection every %s, timeout %s", cfg.Reconnection.CheckInterval, cfg.Reconnection.CheckTimeout)
	ticker := time.NewTicker(cfg.Reconnection.CheckInterval)
	defer ticker.Stop()

	for {
		if s.Session().Closed() {
			break
		}
		select {
		case <-s.done:
			s.log.Infof("done signal trigger")
			return
		case <-ticker.C:
			err := s.Session().Query("SELECT cql_version FROM system.local").Exec()
			if err != nil {
				newSession, err := p.newSession(cfg.Keyspace.Name, cfg.Consistency)
				if err != nil {
					log.Errorf("new session failed: %s", err)
					continue
				}
				s.updateSession(newSession)
			} else {
				continue
			}
		}
	}
}
