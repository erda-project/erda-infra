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

	"github.com/erda-project/erda-infra/pkg/trace"
)

// Future .
type Future interface {
	Get() (interface{}, error)
	Cancel()
}

// Go .
func Go(ctx context.Context, fn func(ctx context.Context) (interface{}, error), opts ...RunOption) Future {
	rOpts := &runOptions{}
	for _, opt := range opts {
		opt(rOpts)
	}
	f := &futureResult{
		result: make(chan *result, 1),
	}
	if rOpts.timeout > 0 {
		f.ctx, f.cancel = context.WithTimeout(ctx, rOpts.timeout)
	} else {
		f.ctx, f.cancel = context.WithCancel(ctx)
	}
	trace.GoWithContext(f.ctx, func(ctx context.Context) {
		defer f.cancel()
		defer close(f.result)
		data, err := fn(ctx)
		f.result <- &result{data, err}
	})
	return f
}

type futureResult struct {
	ctx    context.Context
	cancel func()

	result chan *result
}

type result struct {
	data interface{}
	err  error
}

func (f *futureResult) Get() (interface{}, error) {
	select {
	case <-f.ctx.Done():
		return nil, f.ctx.Err()
	case result := <-f.result:
		return result.data, result.err
	}
}

func (f *futureResult) Cancel() {
	f.cancel()
}
