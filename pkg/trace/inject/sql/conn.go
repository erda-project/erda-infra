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

package sql

import (
	"context"
	"database/sql/driver"

	injectcontext "github.com/erda-project/erda-infra/pkg/trace/inject/context"
)

type wrappedConn struct {
	driver.Conn
}

var (
	_ driver.Pinger             = (*wrappedConn)(nil)
	_ driver.Execer             = (*wrappedConn)(nil) // nolint
	_ driver.ExecerContext      = (*wrappedConn)(nil)
	_ driver.Queryer            = (*wrappedConn)(nil) // nolint
	_ driver.QueryerContext     = (*wrappedConn)(nil)
	_ driver.Conn               = (*wrappedConn)(nil)
	_ driver.ConnPrepareContext = (*wrappedConn)(nil)
	_ driver.ConnBeginTx        = (*wrappedConn)(nil)
	_ driver.SessionResetter    = (*wrappedConn)(nil)
	_ driver.NamedValueChecker  = (*wrappedConn)(nil)
)

func (c *wrappedConn) Ping(ctx context.Context) error {
	pinger, ok := c.Conn.(driver.Pinger)
	if !ok {
		return driver.ErrSkip
	}
	return pinger.Ping(injectcontext.ContextWithSpan(ctx))
}

func (c *wrappedConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	exec, ok := c.Conn.(driver.ExecerContext)
	if ok {
		return exec.ExecContext(injectcontext.ContextWithSpan(context.Background()), query, convertToNamedArgs(args))
	}
	exe, ok := c.Conn.(driver.Execer)
	if ok {
		return exe.Exec(query, args)
	}
	return nil, driver.ErrSkip
}

func (c *wrappedConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execer, ok := c.Conn.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return execer.ExecContext(injectcontext.ContextWithSpan(ctx), query, args)
}

func (c *wrappedConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	queryc, ok := c.Conn.(driver.QueryerContext)
	if ok {
		return queryc.QueryContext(injectcontext.ContextWithSpan(context.Background()), query, convertToNamedArgs(args))
	}
	queryer, ok := c.Conn.(driver.Queryer)
	if ok {
		return queryer.Query(query, args)
	}
	return nil, driver.ErrSkip
}

func (c *wrappedConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := c.Conn.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return queryer.QueryContext(injectcontext.ContextWithSpan(ctx), query, args)
}

func (c *wrappedConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	conn, ok := c.Conn.(driver.ConnPrepareContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return conn.PrepareContext(injectcontext.ContextWithSpan(ctx), query)
}

func (c *wrappedConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	tx, ok := c.Conn.(driver.ConnBeginTx)
	if !ok {
		return nil, driver.ErrSkip
	}
	return tx.BeginTx(injectcontext.ContextWithSpan(ctx), opts)
}

func (c *wrappedConn) ResetSession(ctx context.Context) error {
	r, ok := c.Conn.(driver.SessionResetter)
	if !ok {
		return driver.ErrSkip
	}
	return r.ResetSession(injectcontext.ContextWithSpan(ctx))
}

func (c *wrappedConn) CheckNamedValue(nv *driver.NamedValue) error {
	nvc, ok := c.Conn.(driver.NamedValueChecker)
	if !ok {
		return driver.ErrSkip
	}
	return nvc.CheckNamedValue(nv)
}

func convertToNamedArgs(args []driver.Value) []driver.NamedValue {
	nargs := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		nargs[i] = driver.NamedValue{
			Ordinal: i + 1,
			Value:   arg,
		}
	}
	return nargs
}
