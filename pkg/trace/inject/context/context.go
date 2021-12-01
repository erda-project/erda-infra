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

package hook

import (
	"context"
	"sync"

	"github.com/go-eden/routine"
	"go.opentelemetry.io/otel/trace"
)

const bucketsSize = 128

type (
	contextBucket struct {
		lock sync.RWMutex
		data map[int64]context.Context
	}
	contextBuckets struct {
		buckets [bucketsSize]*contextBucket
	}
)

var goroutineContext contextBuckets

func init() {
	for i := range goroutineContext.buckets {
		goroutineContext.buckets[i] = &contextBucket{
			data: make(map[int64]context.Context),
		}
	}
}

// GetContext .
func GetContext() context.Context {
	goid := routine.Goid()
	idx := goid % bucketsSize
	bucket := goroutineContext.buckets[idx]
	bucket.lock.RLock()
	ctx := bucket.data[goid]
	bucket.lock.RUnlock()
	return ctx
}

// SetContext .
func SetContext(ctx context.Context) {
	goid := routine.Goid()
	idx := goid % bucketsSize
	bucket := goroutineContext.buckets[idx]
	bucket.lock.Lock()
	defer bucket.lock.Unlock()
	bucket.data[goid] = ctx
}

// ClearContext .
func ClearContext() {
	goid := routine.Goid()
	idx := goid % bucketsSize
	bucket := goroutineContext.buckets[idx]
	bucket.lock.Lock()
	defer bucket.lock.Unlock()
	delete(bucket.data, goid)
}

// RunWithContext .
func RunWithContext(ctx context.Context, fn func(ctx context.Context)) {
	SetContext(ctx)
	defer ClearContext()
	fn(ctx)
}

// ContextWithSpan .
func ContextWithSpan(ctx context.Context) context.Context {
	if span := trace.SpanFromContext(ctx); !span.SpanContext().IsValid() {
		pctx := GetContext()
		if pctx != nil {
			if span := trace.SpanFromContext(pctx); span.SpanContext().IsValid() {
				ctx = trace.ContextWithSpan(ctx, span)
			}
		}
	}
	return ctx
}
