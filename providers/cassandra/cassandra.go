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
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	writer "github.com/erda-project/erda-infra/pkg/parallel-writer"
	"github.com/gocql/gocql"
)

// WriterConfig .
type WriterConfig struct {
	Parallelism uint64 `file:"parallelism" default:"4" desc:"parallelism"`
	Batch       struct {
		SizeBytes int           `file:"size_bytes" desc:"cassandra batch failed size bytes"`
		Size      uint64        `file:"size" default:"100" desc:"batch size"`
		Timeout   time.Duration `file:"timeout" default:"30s" desc:"timeout to flush buffer for batch write"`
	} `file:"batch"`
	Retry int `file:"retry" desc:"retry if fail to write"`
}

// SessionConfig .
type SessionConfig struct {
	Keyspace     KeyspaceConfig     `file:"keyspace"`
	Consistency  string             `file:"consistency" default:"LOCAL_ONE"`
	Reconnection ReconnectionConfig `file:"reconnection"`
}

// KeyspaceConfig .
type KeyspaceConfig struct {
	Name        string                    `file:"name" env:"CASSANDRA_KEYSPACE"`
	Auto        bool                      `file:"auto"`
	Replication KeyspaceReplicationConfig `file:"replication"`
}

// KeyspaceReplicationConfig .
type KeyspaceReplicationConfig struct {
	Class  string `file:"class" default:"SimpleStrategy"`
	Factor int32  `file:"factor" default:"2"`
}

// Interface .
type Interface interface {
	CreateKeyspaces(ksc ...*KeyspaceConfig) error
	NewSession(cfg *SessionConfig) (*Session, error)
	NewBatchWriter(session *Session, c *WriterConfig, builderCreator func() StatementBuilder) writer.Writer
}

type config struct {
	Hosts    string        `file:"host" env:"CASSANDRA_ADDR" default:"localhost:9042" desc:"server hosts"`
	Security bool          `file:"security" env:"CASSANDRA_SECURITY_ENABLE" default:"false" desc:"security"`
	Username string        `file:"username" env:"CASSANDRA_SECURITY_USERNAME" default:"" desc:"username"`
	Password string        `file:"password" env:"CASSANDRA_SECURITY_PASSWORD" default:"" desc:"password"`
	Timeout  time.Duration `file:"timeout" env:"CASSANDRA_TIMEOUT" default:"3s" desc:"session timeout"`
}

// provider .
type provider struct {
	Cfg   *config
	Log   logs.Logger
	hosts []string
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.hosts = strings.Split(p.Cfg.Hosts, ",")
	return nil
}

func (p *provider) newSession(keyspace, consistency string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(p.hosts...)
	if p.Cfg.Security && p.Cfg.Username != "" && p.Cfg.Password != "" {
		cluster.Authenticator = &gocql.PasswordAuthenticator{Username: p.Cfg.Username, Password: p.Cfg.Password}
	}
	cluster.Consistency = gocql.ParseConsistency(consistency)
	cluster.Keyspace = keyspace
	cluster.Timeout = p.Cfg.Timeout
	cluster.ConnectTimeout = p.Cfg.Timeout
	return cluster.CreateSession()
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return &service{
		p:   p,
		log: p.Log.Sub(ctx.Caller()),
	}
}

type service struct {
	p    *provider
	log  logs.Logger
	name string
}

func (s *service) CreateKeyspaces(ksc ...*KeyspaceConfig) error {
	var sys *gocql.Session
	defer func() {
		if sys != nil {
			sys.Close()
		}
	}()
	for _, kc := range ksc {
		if sys == nil {
			s, err := s.p.newSession("system", gocql.All.String())
			if err != nil {
				return err
			}
			sys = s
		}
		err := s.createKeySpace(sys, kc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) NewSession(cfg *SessionConfig) (*Session, error) {
	if cfg.Keyspace.Auto {
		err := s.CreateKeyspaces(&cfg.Keyspace)
		if err != nil {
			return nil, err
		}
	}
	session, err := s.p.newSession(cfg.Keyspace.Name, cfg.Consistency)
	if err != nil {
		return nil, fmt.Errorf("new session failed: %w", err)
	}
	ms := &Session{
		session: session,
		log:     s.log.Sub("MySession"),
		done:    make(chan struct{}),
	}

	// workaround for issue: https://github.com/gocql/gocql/issues/831
	// remove it when issue fixed
	if cfg.Reconnection.Enable {
		go ms.checkAndReconnect(s.p, cfg)
	}

	return ms, nil
}

func (s *service) createKeySpace(session *gocql.Session, kc *KeyspaceConfig) error {
	if _, err := session.KeyspaceMetadata(kc.Name); err == nil {
		s.log.Infof("keySpace: %s already existed", kc.Name)
		return nil
	}

	stmt := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH replication={'class':'%s', 'replication_factor':%d}", kc.Name, kc.Replication.Class, kc.Replication.Factor)
	q := session.Query(stmt).Consistency(gocql.All).RetryPolicy(nil)
	defer q.Release()
	s.log.Infof("create keySpace: %s", stmt)
	return q.Exec()
}

func (s *service) NewBatchWriter(session *Session, c *WriterConfig, builderCreator func() StatementBuilder) writer.Writer {
	return writer.ParallelBatch(func(uint64) writer.Writer {
		return &batchWriter{
			session:        session,
			builder:        builderCreator(),
			retry:          c.Retry,
			retryDuration:  3 * time.Second,
			log:            s.log,
			batchSizeBytes: c.Batch.SizeBytes,
		}
	}, c.Parallelism, c.Batch.Size, c.Batch.Timeout, s.batchWriteError)
}

func (s *service) batchWriteError(err error) error {
	s.log.Errorf("fail to write cassandra: %s", err)
	return nil // skip error
}

func init() {
	servicehub.Register("cassandra", &servicehub.Spec{
		Services:    []string{"cassandra"},
		Description: "cassandra",
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
		},
		ConfigFunc: func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
