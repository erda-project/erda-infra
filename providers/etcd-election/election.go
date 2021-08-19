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

package election

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/recallsong/go-utils/errorx"
	uuid "github.com/satori/go.uuid"
)

// Node .
type Node struct {
	ID string
}

// Action .
type Action int32

// Action values
const (
	ActionPut = Action(iota + 1)
	ActionDelete
)

func (a Action) String() string {
	switch a {
	case ActionPut:
		return "put"
	case ActionDelete:
		return "delete"
	}
	return "unknown"
}

// Event .
type Event struct {
	Action Action
	Node   Node
}

// WatchOption .
type WatchOption interface{}

// Interface .
type Interface interface {
	Node() Node
	Nodes() ([]Node, error)
	Leader() (*Node, error)
	IsLeader() bool
	ResignLeader() error
	OnLeader(handler func(context.Context))
	Watch(ctx context.Context, opts ...WatchOption) <-chan Event
}

type config struct {
	Prefix string `file:"root_path" default:"etcd-election"`
	NodeID string `file:"node_id"`
}

type provider struct {
	Cfg    *config
	Log    logs.Logger
	Client *clientv3.Client `autowired:"etcd-client"`
	prefix string

	lock           sync.RWMutex
	leaderHandlers []func(ctx context.Context)
	cancelHandler  func()
	election       *concurrency.Election
	session        *concurrency.Session
	iAmLeader      bool
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.Cfg.Prefix = filepath.Clean("/" + p.Cfg.Prefix)
	p.prefix = p.Cfg.Prefix + "/"
	if len(p.Cfg.NodeID) <= 0 {
		p.Cfg.NodeID = uuid.NewV4().String()
	}
	p.Log.Info("my node id: ", p.Cfg.NodeID)
	return nil
}

func (p *provider) reset(session *concurrency.Session) {
	session.Close()
	p.lock.Lock()
	p.session, p.election = nil, nil
	p.iAmLeader = false
	p.lock.Unlock()
}

func (p *provider) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		session, err := p.newSession(ctx, 5*time.Second)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			p.Log.Errorf("fail to NewSession: %s", err)
			time.Sleep(2 * time.Second)
			continue
		}

		election := concurrency.NewElection(session, p.Cfg.Prefix)
		p.lock.Lock()
		p.session, p.election = session, election
		p.lock.Unlock()
		if err = election.Campaign(ctx, p.Cfg.NodeID); err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			p.reset(session)
			p.Log.Errorf("fail to Campaign: %s", err, reflect.TypeOf(err))
			time.Sleep(1 * time.Second)
			continue
		}

		// Let's say A is leader and B is non-leader.
		// The etcd server's stopped and it's restarted after a while like 10 seconds.
		// The campaign of B exited with nil after connection was restored.
		select {
		case <-session.Done():
			p.reset(session)
			continue
		default:
		}

		p.Log.Infof("I am leader ! Node is %q", p.Cfg.NodeID)

		p.runHandlers()
	loop:
		for {
			select {
			case <-session.Done():
				p.resignLeader()
				break loop
			case <-ctx.Done():
				p.resignLeader()
			}
		}
	}
}

func (p *provider) newSession(ctx context.Context, ttl time.Duration) (*concurrency.Session, error) {
	opts := []concurrency.SessionOption{concurrency.WithContext(ctx)}
	seconds := int(ttl.Seconds())
	if seconds > 0 {
		opts = append(opts, concurrency.WithTTL(seconds))
	}
	return concurrency.NewSession(p.Client, opts...)
}

func (p *provider) runHandlers() {
	ctx, cancel := context.WithCancel(context.Background())
	p.lock.Lock()
	p.iAmLeader = true
	p.cancelHandler = cancel
	p.lock.Unlock()
	for _, h := range p.leaderHandlers {
		go func(h func(context.Context)) {
			h(ctx)
		}(h)
	}
}

func (p *provider) Node() Node {
	return Node{ID: p.Cfg.NodeID}
}

func (p *provider) Nodes() ([]Node, error) {
	resp, err := p.Client.Get(context.Background(), p.prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var nodes []Node
	for _, kv := range resp.Kvs {
		nodes = append(nodes, Node{ID: string(kv.Value)})
	}
	return nodes, nil
}

func (p *provider) Leader() (*Node, error) {
	if p.IsLeader() {
		node := p.Node()
		return &node, nil
	}
	clientv3.WithFirstCreate()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := p.Client.Get(ctx, p.prefix, clientv3.WithFirstCreate()...)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	node := &Node{ID: string(resp.Kvs[0].Value)}
	return node, nil
}

func (p *provider) IsLeader() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.iAmLeader
}

func (p *provider) ResignLeader() error {
	err := p.resignLeader()
	if err != nil {
		p.Log.Warnf("fail to resign leader: %s", err)
	}
	return err
}

func (p *provider) resignLeader() error {
	var election *concurrency.Election
	var session *concurrency.Session

	p.lock.Lock()
	if !p.iAmLeader {
		p.lock.Unlock()
		return nil
	}
	p.iAmLeader = false
	p.cancelHandler()
	p.cancelHandler = nil
	election = p.election
	session = p.session
	p.session, p.election = nil, nil
	p.lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	var errs errorx.Errors
	err := election.Resign(ctx)
	if err != nil {
		errs.Append(err)
	}
	err = session.Close()
	if err != nil {
		errs.Append(err)
	}
	return errs.MaybeUnwrap()
}

func (p *provider) OnLeader(handler func(context.Context)) {
	p.leaderHandlers = append(p.leaderHandlers, handler)
}

func (p *provider) Watch(ctx context.Context, opts ...WatchOption) <-chan Event {
	notify := make(chan Event, 8)
	go func() {
		defer func() {
			close(notify)
			p.Log.Debug("election watcher exited")
		}()
		opts := []clientv3.OpOption{clientv3.WithPrefix()}
		for func() bool {
			wctx, wcancel := context.WithCancel(context.Background())
			defer wcancel()
			wch := p.Client.Watch(wctx, p.prefix, opts...)
			for {
				select {
				case wr, ok := <-wch:
					if !ok {
						return true
					} else if wr.Err() != nil {
						p.Log.Errorf("election watcher error: %s", wr.Err())
						return true
					}
					for _, ev := range wr.Events {
						if ev.Kv == nil {
							continue
						}
						switch ev.Type {
						case mvccpb.PUT:
							notify <- Event{
								Action: ActionDelete,
								Node:   Node{ID: string(ev.Kv.Value)},
							}
						case mvccpb.DELETE:
							notify <- Event{
								Action: ActionDelete,
								Node:   Node{ID: string(ev.Kv.Value)},
							}
						}
					}
				case <-ctx.Done():
					return false
				}
			}
		}() {
		}
	}()
	return notify
}

func init() {
	servicehub.Register("etcd-election", &servicehub.Spec{
		Services: []string{"etcd-election"},
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
		},
		Dependencies: []string{"etcd"},
		ConfigFunc:   func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
