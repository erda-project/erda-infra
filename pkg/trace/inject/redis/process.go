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

package redis

import (
	"github.com/go-redis/redis"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	injectcontext "github.com/erda-project/erda-infra/pkg/trace/inject/context"
)

func newProcessWrapper(cfg *config) func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
	return func(process func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) (err error) {
			ctx := injectcontext.GetContext()
			if ctx != nil {
				record := cfg.SpanOptions.AllowRoot || trace.SpanContextFromContext(ctx).IsValid()
				if record && !cfg.SpanOptions.Ping && cmd.Name() == "ping" {
					record = false
				}
				if record {
					var span trace.Span
					ctx, span = cfg.Tracer.Start(ctx, cfg.SpanNameFormatter.Format(ctx, cmd),
						trace.WithSpanKind(trace.SpanKindClient),
						trace.WithAttributes(withDBStatement(cfg, cmd)...),
					)
					defer func() {
						if err != nil {
							recordSpanError(span, cfg.SpanOptions, err)
						} else {
							span.SetStatus(codes.Ok, "")
						}
						span.End()
					}()
				}
			}
			err = process(cmd)
			return err
		}
	}
}

func newProcessPipeline(cfg *config) func(process func([]redis.Cmder) error) func([]redis.Cmder) error {
	return func(process func(cmds []redis.Cmder) error) func([]redis.Cmder) error {
		return func(cmds []redis.Cmder) (err error) {
			ctx := injectcontext.GetContext()
			if ctx != nil {
				if cfg.SpanOptions.AllowRoot || trace.SpanContextFromContext(ctx).IsValid() {
					var span trace.Span
					ctx, span = cfg.Tracer.Start(ctx, cfg.SpanNameFormatter.FormatBatch(ctx, cmds),
						trace.WithSpanKind(trace.SpanKindClient),
						trace.WithAttributes(withBatchDBStatement(cfg, cmds)...),
					)
					defer func() {
						if err != nil {
							recordSpanError(span, cfg.SpanOptions, err)
						} else {
							span.SetStatus(codes.Ok, "")
						}
						span.End()
					}()
				}
			}
			err = process(cmds)
			return err
		}
	}
}

func recordSpanError(span trace.Span, opts SpanOptions, err error) {
	if span == nil {
		return
	}
	if opts.RecordError != nil && !opts.RecordError(err) {
		return
	}
	switch err {
	case nil:
		return
	default:
		span.RecordError(err)
		span.SetStatus(codes.Error, "")
	}
}
