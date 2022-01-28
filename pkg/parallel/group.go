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

package parallel

import (
	"context"
	"sync"
	"time"

	"github.com/erda-project/erda-infra/pkg/trace"
)

// Group .
type Group struct {
	ctx    context.Context
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// GoGroup
func GoGroup(ctx context.Context, opts ...RunOption) *Group {
	rOpts := &runOptions{}
	for _, opt := range opts {
		opt(rOpts)
	}
	g := &Group{}
	if rOpts.timeout > 0 {
		g.ctx, g.cancel = context.WithTimeout(ctx, rOpts.timeout)
	} else {
		g.ctx, g.cancel = context.WithCancel(ctx)
	}
	return g
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(fn func(ctx context.Context) error) {
	g.wg.Add(1)
	trace.GoWithContext(g.ctx, func(ctx context.Context) {
		defer g.wg.Done()
		if err := fn(ctx); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	})
}

// RunOption .
type RunOption func(opts *runOptions)

type runOptions struct {
	timeout time.Duration
}

// WithTimeout .
func WithTimeout(timeout time.Duration) RunOption {
	return func(opts *runOptions) {
		opts.timeout = timeout
	}
}
