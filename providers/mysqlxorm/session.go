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
	"github.com/xormplus/xorm"
)

type (
	session struct {
		*xorm.Session
		needAutoClose bool
	}
	SessionOption func(s *session)
)

func WithSession(passedInSession *session) SessionOption {
	return func(s *session) {
		s.Session = passedInSession.Session
	}
}

func (p *provider) NewSession(opts ...SessionOption) *session {
	tx := &session{}

	for _, opt := range opts {
		opt(tx)
	}

	// set default session
	if tx.Session == nil {
		tx.Session = p.db.NewSession()
		tx.needAutoClose = true
	}

	return tx
}

func (tx *session) Close() {
	if tx.needAutoClose {
		tx.Session.Close()
	}
	return
}
