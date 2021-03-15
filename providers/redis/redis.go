// Author: recallsong
// Email: songruiguo@qq.com

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

func (d *define) Service() []string { return []string{"redis", "redis-client"} }
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
	C       *config
	L       logs.Logger
	client  *redis.Client
	clients map[int]*redis.Client
	lock    sync.Mutex
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	if p.C.DB <= 0 {
		return nil
	}
	c, err := p.Open(p.C.DB)
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
	c, _ := p.Open(p.C.DB)
	return c
}

func (p *provider) Open(db int) (*redis.Client, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if c, ok := p.clients[db]; ok {
		return c, nil
	}
	var c *redis.Client
	if p.C.MasterName != "" && p.C.SentinelsAddr != "" {
		addrs := strings.Split(p.C.SentinelsAddr, ",")
		c = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:         p.C.MasterName,
			SentinelAddrs:      addrs,
			Password:           p.C.Password,
			DB:                 db,
			MaxRetries:         p.C.MaxRetries,
			DialTimeout:        p.C.DialTimeout,
			ReadTimeout:        p.C.ReadTimeout,
			WriteTimeout:       p.C.WriteTimeout,
			PoolSize:           p.C.PoolSize,
			PoolTimeout:        p.C.PoolTimeout,
			IdleTimeout:        p.C.IdleTimeout,
			IdleCheckFrequency: p.C.IdleCheckFrequency,
		})
	} else if p.C.Addr != "" {
		c = redis.NewClient(&redis.Options{
			Addr:               p.C.Addr,
			Password:           p.C.Password,
			DB:                 db,
			MaxRetries:         p.C.MaxRetries,
			DialTimeout:        p.C.DialTimeout,
			ReadTimeout:        p.C.ReadTimeout,
			WriteTimeout:       p.C.WriteTimeout,
			PoolSize:           p.C.PoolSize,
			PoolTimeout:        p.C.PoolTimeout,
			IdleTimeout:        p.C.IdleTimeout,
			IdleCheckFrequency: p.C.IdleCheckFrequency,
		})
	} else {
		err := fmt.Errorf("redis config error: no addr or sentinel")
		p.L.Error(err)
		return nil, err
	}

	if pong, err := c.Ping().Result(); err != nil {
		p.L.Errorf("redis ping error: %s", err)
		return nil, err
	} else if pong != "PONG" {
		err := fmt.Errorf("redis ping result '%s' is not PONG", pong)
		p.L.Error(err)
		return nil, err
	} else {
		p.L.Infof("open redis db %d and ping ok", db)
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
