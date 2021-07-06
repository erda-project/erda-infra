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

package mutex

import (
	"context"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/etcd"
)

// Mutex .
type Mutex interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
	Close() error
}

// Interface .
type Interface interface {
	NewWithTTL(ctx context.Context, key string, ttl time.Duration) (Mutex, error)
	New(ctx context.Context, key string) (Mutex, error)
}

var mutexType = reflect.TypeOf((*Mutex)(nil)).Elem()

type config struct {
	RootPath   string `file:"root_path"`
	DefaultKey string `file:"default_key"`
}

type provider struct {
	Cfg         *config
	Log         logs.Logger
	etcd        etcd.Interface
	instances   map[string]Mutex
	inProcMutex *inProcMutex
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.etcd = ctx.Service("etcd").(etcd.Interface)
	p.Cfg.RootPath = filepath.Clean("/" + p.Cfg.RootPath)
	return nil
}

// NewWithTTL .
func (p *provider) NewWithTTL(ctx context.Context, key string, ttl time.Duration) (Mutex, error) {
	opts := []concurrency.SessionOption{concurrency.WithContext(ctx)}
	seconds := int(ttl.Seconds())
	if seconds > 0 {
		opts = append(opts, concurrency.WithTTL(seconds))
	}
	session, err := concurrency.NewSession(p.etcd.Client(), opts...)
	if err != nil {
		return nil, err
	}
	key = filepath.Clean(filepath.Join(p.Cfg.RootPath, key))
	mutex := concurrency.NewMutex(session, key)
	return &etcdMutex{
		log: p.Log,
		key: key,
		s:   session,
		mu:  mutex,
	}, nil
}

func (p *provider) New(ctx context.Context, key string) (Mutex, error) {
	return p.NewWithTTL(ctx, key, time.Duration(0))
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Type() == mutexType {
		key := ctx.Tags().Get("mutex-key")
		if len(key) < 0 {
			key = p.Cfg.DefaultKey
		}
		if len(key) < 0 {
			return p.inProcMutex
		}
		m, err := p.New(context.Background(), key)
		if err != nil {
			p.Log.Errorf("fail to create mutex for key: %q", key)
		}
		return m
	}
	return p
}

type etcdMutex struct {
	log  logs.Logger
	key  string
	s    *concurrency.Session
	mu   *concurrency.Mutex
	lock sync.Mutex
}

func (m *etcdMutex) Lock(ctx context.Context) error {
	m.lock.Lock()
	err := m.mu.Lock(ctx)
	if err == nil {
		m.log.Debugf("locked key %q", m.key)
	} else {
		m.log.Errorf("fail to lock key %q", m.key)
	}
	return err
}

func (m *etcdMutex) Unlock(ctx context.Context) error {
	m.lock.Unlock()
	err := m.mu.Unlock(ctx)
	if err == nil {
		m.log.Debugf("unlocked key %q", m.key)
	} else {
		m.log.Errorf("fail to unlock key %q", m.key)
	}
	return err
}

func (m *etcdMutex) Close() error { return m.s.Close() }

type inProcMutex struct {
	lock sync.Mutex
}

func (m *inProcMutex) Lock(ctx context.Context) error {
	m.lock.Lock()
	return nil
}

func (m *inProcMutex) Unlock(ctx context.Context) error {
	m.lock.Unlock()
	return nil
}

func (m *inProcMutex) Close() error { return nil }

func init() {
	servicehub.Register("etcd-mutex", &servicehub.Spec{
		Services: []string{"etcd-mutex"},
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
			mutexType,
		},
		Dependencies: []string{"etcd"},
		Description:  "distributed lock implemented by etcd",
		ConfigFunc:   func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{
				instances:   make(map[string]Mutex),
				inProcMutex: &inProcMutex{},
			}
		},
	})
}
