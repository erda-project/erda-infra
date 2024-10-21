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
	"errors"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

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

// ErrClosed mutex already closed
var ErrClosed = errors.New("mutex closed")

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
	lock        *sync.Mutex
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.etcd = ctx.Service("etcd").(etcd.Interface)
	p.Cfg.RootPath = filepath.Clean("/" + p.Cfg.RootPath)
	return nil
}

// NewWithTTL .
func (p *provider) NewWithTTL(ctx context.Context, key string, ttl time.Duration) (Mutex, error) {
	ctx, cancel := context.WithCancel(ctx)
	opts := []concurrency.SessionOption{concurrency.WithContext(ctx)}
	seconds := int(ttl.Seconds())
	if seconds > 0 {
		opts = append(opts, concurrency.WithTTL(seconds))
	}
	mutex := &etcdMutex{
		log:        p.Log,
		key:        filepath.Clean(filepath.Join(p.Cfg.RootPath, key)),
		client:     p.etcd.Client(),
		opts:       opts,
		inProcLock: make(chan struct{}, 1),
		ctx:        ctx,
		cancel:     cancel,
	}
	p.lock.Lock()
	defer p.lock.Unlock()
	p.instances[key] = mutex
	return mutex, nil
}

// New .
func (p *provider) New(ctx context.Context, key string) (Mutex, error) {
	p.lock.Lock()
	if ins, ok := p.instances[key]; ok {
		return ins, nil
	}
	p.lock.Unlock()
	return p.NewWithTTL(ctx, key, time.Duration(0))
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Type() == mutexType {
		key := ctx.Tags().Get("mutex-key")
		if len(key) <= 0 {
			key = p.Cfg.DefaultKey
		}
		if len(key) <= 0 {
			p.Log.Debugf("in-proc mutex for provider %q", ctx.Caller())
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
	log    logs.Logger
	key    string
	client *clientv3.Client
	opts   []concurrency.SessionOption
	ctx    context.Context
	cancel context.CancelFunc

	lock       sync.Mutex
	s          *concurrency.Session
	mu         *concurrency.Mutex
	inProcLock chan struct{}
}

func (m *etcdMutex) resetSession() (*concurrency.Mutex, error) {
	m.close()
	s, mu, err := m.newSession()
	if err != nil {
		return nil, err
	}
	m.s, m.mu = s, mu
	return mu, nil
}

func (m *etcdMutex) newSession() (*concurrency.Session, *concurrency.Mutex, error) {
	session, err := concurrency.NewSession(m.client, m.opts...)
	if err != nil {
		m.log.Debugf("failed to new session for key %q: %s", m.key, err)
		return nil, nil, err
	}
	m.log.Debugf("new session for key %q", m.key)
	return session, concurrency.NewMutex(session, m.key), nil
}

func (m *etcdMutex) Lock(ctx context.Context) (err error) {
	d := time.Second
	sleep := func() bool {
		select {
		case <-time.After(d):
		case <-ctx.Done():
			return false
		}
		if d < 8*time.Second {
			d = d * 2
		}
		return true
	}

	select {
	case m.inProcLock <- struct{}{}:
	case <-m.ctx.Done():
		return ErrClosed
	case <-ctx.Done():
		return context.Canceled
	}

	for {
		m.lock.Lock()
		select {
		case <-m.ctx.Done():
			m.lock.Unlock()
			return ErrClosed
		case <-ctx.Done():
			m.lock.Unlock()
			return context.Canceled
		default:
		}
		mu := m.mu
		if err != nil || mu == nil {
			mu, err = m.resetSession()
			if err != nil {
				m.lock.Unlock()
				if errors.Is(err, context.Canceled) {
					return err
				}
				sleep()
				continue
			}
		}
		m.lock.Unlock()

		err = mu.Lock(ctx)
		if err != nil {
			m.log.Errorf("failed to lock key %q: %s", m.key, err)
			if errors.Is(err, context.Canceled) {
				return err
			}
			continue
		}
		m.log.Debugf("locked key %q", m.key)
		return nil
	}
}

func (m *etcdMutex) Unlock(ctx context.Context) (err error) {
	select {
	case <-m.inProcLock:
	case <-m.ctx.Done():
		return ErrClosed
	case <-ctx.Done():
		return context.Canceled
	}

	m.lock.Lock()
	mu := m.mu
	if mu != nil {
		err = m.mu.Unlock(ctx)
	}
	m.lock.Unlock()

	if err != nil {
		m.log.Errorf("failed to unlock key %q: %s", m.key, err)
		return err
	}
	m.log.Debugf("unlocked key %q", m.key)
	return err
}

func (m *etcdMutex) Close() error {
	m.lock.Lock()
	select {
	case <-m.ctx.Done():
		m.lock.Unlock()
		return nil
	default:
		m.cancel()
	}
	err := m.close()
	if errors.Is(err, context.Canceled) {
		err = nil
	}
	m.lock.Unlock()
	return err
}

func (m *etcdMutex) close() (err error) {
	if m.s != nil {
		err = m.s.Close()
		m.s, m.mu = nil, nil
	}
	return err
}

type inProcMutex struct {
	lock chan struct{}
}

func (m *inProcMutex) Lock(ctx context.Context) error {
	select {
	case m.lock <- struct{}{}:
	case <-ctx.Done():
		return context.Canceled
	}
	return nil
}

func (m *inProcMutex) Unlock(ctx context.Context) error {
	select {
	case <-m.lock:
	case <-ctx.Done():
		return context.Canceled
	}
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
				inProcMutex: &inProcMutex{lock: make(chan struct{}, 1)},
				lock:        &sync.Mutex{},
			}
		},
	})
}
