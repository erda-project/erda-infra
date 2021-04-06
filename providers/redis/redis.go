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

package redis

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/go-redis/redis"
)

// Interface .
type Interface interface {
	DB() *redis.Client
	Open(db int) (*redis.Client, error)
}

var (
	interfaceType = reflect.TypeOf((*Interface)(nil)).Elem()
	clientType    = reflect.TypeOf((*redis.Client)(nil))
)

type config struct {
	Addr          string `file:"addr" env:"REDIS_ADDR"`
	MasterName    string `file:"master_name" env:"REDIS_MASTER_NAME"`
	SentinelsAddr string `file:"sentinels_addr" env:"REDIS_SENTINELS_ADDR"`
	Password      string `file:"password" env:"REDIS_PASSWORD"`
	DB            int    `file:"db" env:"REDIS_DB"`

	MaxRetries int `file:"max_retries" env:"REDIS_MAX_RETRIES"`

	DialTimeout  time.Duration `file:"dial_timeout" env:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `file:"read_timeout" env:"REDIS_READ_TIMEOUT"`
	WriteTimeout time.Duration `file:"write_timeout" env:"REDIS_WRITE_TIMEOUT"`

	PoolSize           int           `file:"pool_size" env:"REDIS_POOL_SIZE"`
	PoolTimeout        time.Duration `file:"pool_timeout" env:"REDIS_POOL_TIMEOUT"`
	IdleTimeout        time.Duration `file:"idle_timeout" env:"REDIS_IDLE_TIMEOUT"`
	IdleCheckFrequency time.Duration `file:"idle_check_frequency" env:"REDIS_IDLE_CHECK_FREQUENCY"`
}

type define struct{}

func (d *define) Services() []string { return []string{"redis", "redis-client"} }
func (d *define) Types() []reflect.Type {
	return []reflect.Type{
		interfaceType, clientType,
	}
}
func (d *define) Summary() string     { return "redis" }
func (d *define) Description() string { return d.Summary() }
func (d *define) Config() interface{} { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{clients: make(map[int]*redis.Client)}
	}
}

// provider .
type provider struct {
	Cfg     *config
	Log     logs.Logger
	client  *redis.Client
	clients map[int]*redis.Client
	lock    sync.Mutex
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	if p.Cfg.DB <= 0 {
		return nil
	}
	c, err := p.Open(p.Cfg.DB)
	if err != nil {
		return err
	}
	p.client = c
	return nil
}

func (p *provider) DB() *redis.Client {
	if p.client != nil {
		return p.client
	}
	c, _ := p.Open(p.Cfg.DB)
	return c
}

func (p *provider) Open(db int) (*redis.Client, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if c, ok := p.clients[db]; ok {
		return c, nil
	}
	var c *redis.Client
	if p.Cfg.MasterName != "" && p.Cfg.SentinelsAddr != "" {
		addrs := strings.Split(p.Cfg.SentinelsAddr, ",")
		c = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:         p.Cfg.MasterName,
			SentinelAddrs:      addrs,
			Password:           p.Cfg.Password,
			DB:                 db,
			MaxRetries:         p.Cfg.MaxRetries,
			DialTimeout:        p.Cfg.DialTimeout,
			ReadTimeout:        p.Cfg.ReadTimeout,
			WriteTimeout:       p.Cfg.WriteTimeout,
			PoolSize:           p.Cfg.PoolSize,
			PoolTimeout:        p.Cfg.PoolTimeout,
			IdleTimeout:        p.Cfg.IdleTimeout,
			IdleCheckFrequency: p.Cfg.IdleCheckFrequency,
		})
	} else if p.Cfg.Addr != "" {
		c = redis.NewClient(&redis.Options{
			Addr:               p.Cfg.Addr,
			Password:           p.Cfg.Password,
			DB:                 db,
			MaxRetries:         p.Cfg.MaxRetries,
			DialTimeout:        p.Cfg.DialTimeout,
			ReadTimeout:        p.Cfg.ReadTimeout,
			WriteTimeout:       p.Cfg.WriteTimeout,
			PoolSize:           p.Cfg.PoolSize,
			PoolTimeout:        p.Cfg.PoolTimeout,
			IdleTimeout:        p.Cfg.IdleTimeout,
			IdleCheckFrequency: p.Cfg.IdleCheckFrequency,
		})
	} else {
		err := fmt.Errorf("redis config error: no addr or sentinel")
		p.Log.Error(err)
		return nil, err
	}

	if pong, err := c.Ping().Result(); err != nil {
		p.Log.Errorf("redis ping error: %s", err)
		return nil, err
	} else if pong != "PONG" {
		err := fmt.Errorf("redis ping result '%s' is not PONG", pong)
		p.Log.Error(err)
		return nil, err
	} else {
		p.Log.Infof("open redis db %d and ping ok", db)
	}

	p.clients[db] = c
	return c, nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Service() == "redis-client" || ctx.Type() == clientType {
		return p.DB()
	}
	return p
}

func init() {
	servicehub.RegisterProvider("redis", &define{})
}
