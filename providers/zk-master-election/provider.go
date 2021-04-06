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

package election

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/zookeeper"
	"github.com/go-zookeeper/zk"
)

// Event .
type Event interface {
	IsConnected() bool
	IsMaster() bool
}

// Listener .
type Listener func(Event)

// Interface .
type Interface interface {
	IsConnected() bool
	IsMaster() bool
	Watch(Listener)
}

type define struct{}

func (d *define) Services() []string     { return []string{"zk-master-election"} }
func (d *define) Dependencies() []string { return []string{"zookeeper"} }
func (d *define) Types() []reflect.Type {
	return []reflect.Type{reflect.TypeOf((*Interface)(nil)).Elem()}
}
func (d *define) Description() string { return "master election implemented by zookeeper" }
func (d *define) Config() interface{} { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{
			closeCh:  make(chan struct{}),
			watchers: make(map[string][]Listener),
		}
	}
}

type config struct {
	RootPath   string `file:"root_path"`
	MasterNode string `file:"master_node" default:"master-node-key"`
	masterPath string
}

type provider struct {
	Cfg     *config
	Log     logs.Logger
	zk      zookeeper.Interface
	closeCh chan struct{}

	isConnected int32
	isMaster    int32
	keys        []string
	watchers    map[string][]Listener
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.zk = ctx.Service("zookeeper").(zookeeper.Interface)
	p.Cfg.RootPath = filepath.Clean("/" + p.Cfg.RootPath)
	p.Cfg.MasterNode = filepath.Clean(p.Cfg.MasterNode)
	p.Cfg.masterPath = filepath.Join(p.Cfg.RootPath, p.Cfg.MasterNode)
	return nil
}

func (p *provider) run() error {
	for {
		conn, ch, err := p.zk.Connect()
		if err != nil {
			p.Log.Errorf("fail to connect zookeeper: %s", err)
			select {
			case <-p.closeCh:
				if conn != nil {
					conn.Close()
				}
				return err
			default:
				time.Sleep(3 * time.Second)
			}
			continue
		}
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())
		timer := time.After(p.zk.SessionTimeout())
		for {
			var exit bool
			select {
			case event := <-ch:
				if event.Type != zk.EventSession {
					continue
				}
				switch event.State {
				case zk.StateConnected:
					atomic.StoreInt32(&p.isConnected, 1)
					p.Log.Info("connected to zookeeper successfully")
					err := p.election(conn)
					if err != nil {
						break
					}
					wg.Add(1)
					go p.watchMasterNode(ctx, &wg, conn)
					continue
				case zk.StateConnectedReadOnly, zk.StateConnecting, zk.StateHasSession, zk.StateSaslAuthenticated, zk.StateUnknown:
					continue
				case zk.StateExpired, zk.StateAuthFailed, zk.StateDisconnected:
					break
				default:
					p.Log.Errorf("unknown event: %v", event)
					continue
				}
			case <-timer:
				if !p.IsConnected() {
					p.Log.Errorf("connect to zookeeper timeout")
					break
				}
				continue
			case <-p.closeCh:
				exit = true
			}
			cancel()
			atomic.StoreInt32(&p.isMaster, 0)
			atomic.StoreInt32(&p.isConnected, 0)
			wg.Wait()
			conn.Close()
			p.Log.Info("disconnected zookeeper")
			if exit {
				return nil
			}
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func (p *provider) IsConnected() bool {
	return atomic.LoadInt32(&p.isConnected) != 0
}

func (p *provider) IsMaster() bool {
	return atomic.LoadInt32(&p.isMaster) != 0
}

func (p *provider) makePath(conn *zk.Conn, path string) error {
	exist, _, err := conn.Exists(path)
	if err != nil {
		return err
	}
	if !exist {
		createdPath, err := conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return fmt.Errorf("fail to create path %q: %s", path, err)
		}
		if path != createdPath {
			return fmt.Errorf("create different path %q != %q", createdPath, path)
		}
		p.Log.Infof("created path %q", path)
	}
	return nil
}

type stateEvent struct {
	isConnected bool
	isMaster    bool
}

func (c *stateEvent) IsConnected() bool { return c.isConnected }
func (c *stateEvent) IsMaster() bool    { return c.isMaster }

func (p *provider) election(conn *zk.Conn) error {
	err := p.makePath(conn, p.Cfg.RootPath)
	if err != nil {
		return err
	}
	createdPath, err := conn.Create(p.Cfg.masterPath, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		if !strings.Contains(err.Error(), "exists") {
			err = fmt.Errorf("fail to create path %q: %s", p.Cfg.masterPath, err)
			p.Log.Error(err)
			return err
		}
	} else if createdPath != p.Cfg.masterPath {
		err = fmt.Errorf("create different path %q != %q", createdPath, p.Cfg.masterPath)
		p.Log.Error(err)
		return err
	}
	isMaster := err == nil
	if isMaster {
		atomic.StoreInt32(&p.isMaster, 0)
		p.Log.Infof("election finish, i am slave")
	} else {
		atomic.StoreInt32(&p.isMaster, 1)
		p.Log.Infof("election success, i am master")
	}
	ctx := &stateEvent{
		isMaster:    isMaster,
		isConnected: p.IsConnected(),
	}
	for _, key := range p.keys {
		for _, w := range p.watchers[key] {
			w(ctx)
		}
	}
	return nil
}

func (p *provider) watchMasterNode(ctx context.Context, wg *sync.WaitGroup, conn *zk.Conn) {
	defer wg.Done()
loop:
	for {
		_, _, ch, err := conn.ChildrenW(p.Cfg.masterPath)
		if err != nil {
			p.Log.Errorf("fail to watch path %q: %s", p.Cfg.masterPath, err)
			select {
			case <-ctx.Done():
			default:
				time.Sleep(3 * time.Second)
			}
			continue
		}
		p.Log.Infof("start watch path %q", p.Cfg.masterPath)
		defer p.Log.Infof("exit waith path %q", p.Cfg.masterPath)
		for {
			select {
			case event, ok := <-ch:
				if !ok {
					continue loop
				}
				if event.Type == zk.EventNodeDeleted {
					err := p.election(conn)
					if err != nil {
						continue loop
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (p *provider) Start() error {
	return p.run()
}

func (p *provider) Close() error {
	close(p.closeCh)
	return nil
}

type service struct {
	p    *provider
	name string
}

func (s *service) IsMaster() bool {
	return s.p.IsMaster()
}

func (s *service) IsConnected() bool {
	return s.p.IsConnected()
}

func (s *service) Watch(ln Listener) {
	list, ok := s.p.watchers[s.name]
	if !ok {
		s.p.keys = append(s.p.keys, s.name)
	}
	s.p.watchers[s.name] = append(list, ln)
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return &service{
		p:    p,
		name: ctx.Caller(),
	}
}

func init() {
	servicehub.RegisterProvider("zk-master-election", &define{})
}
