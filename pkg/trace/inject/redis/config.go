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
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/erda-project/erda-infra/pkg/trace/inject/redis"
)

// SpanNameFormatter is an interface that used to format span names.
type SpanNameFormatter interface {
	Format(ctx context.Context, cmd redis.Cmder) string
	FormatBatch(ctx context.Context, cmds []redis.Cmder) string
}

type config struct {
	TracerProvider trace.TracerProvider
	Tracer         trace.Tracer

	SpanOptions SpanOptions

	DBSystem string

	// Attributes will be set to each span.
	Attributes []attribute.KeyValue

	// SpanNameFormatter will be called to produce span's name.
	// Default use method as span name
	SpanNameFormatter SpanNameFormatter
}

// SpanOptions holds configuration of tracing span to decide
// whether to enable some features.
// by default all options are set to false intentionally when creating a wrapped
// driver and provide the most sensible default with both performance and
// security in mind.
type SpanOptions struct {
	// Ping, if set to true, will enable the creation of spans on Ping requests.
	Ping bool

	// DisableStatement if set to true, will suppress db.statement in spans.
	DisableStatement bool

	// RecordError, if set, will be invoked with the current error, and if the func returns true
	// the record will be recorded on the current span.
	RecordError func(err error) bool

	// AllowRoot, if set to true, will create root spans in absence of existing spans or even context.
	AllowRoot bool
}

type defaultSpanNameFormatter struct{}

func (f *defaultSpanNameFormatter) Format(ctx context.Context, cmd redis.Cmder) string {
	return cmd.Name()
}

func (f *defaultSpanNameFormatter) FormatBatch(ctx context.Context, cmds []redis.Cmder) string {
	return "batch"
}

// newConfig returns a config with all Options set.
func newConfig(dbSystem string, options ...Option) *config {
	cfg := config{
		TracerProvider:    otel.GetTracerProvider(),
		DBSystem:          dbSystem,
		SpanNameFormatter: &defaultSpanNameFormatter{},
		SpanOptions:       SpanOptions{},
	}
	for _, opt := range options {
		opt.Apply(&cfg)
	}

	if cfg.DBSystem != "" {
		cfg.Attributes = append(cfg.Attributes,
			semconv.DBSystemKey.String(cfg.DBSystem),
		)
		cfg.Attributes = cfg.Attributes[:len(cfg.Attributes):len(cfg.Attributes)]
	}
	cfg.Tracer = cfg.TracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion(Version()),
	)
	return &cfg
}

func withDBStatement(cfg *config, cmd redis.Cmder) []attribute.KeyValue {
	if cfg.SpanOptions.DisableStatement {
		return cfg.Attributes
	}
	return append(cfg.Attributes, semconv.DBStatementKey.String(getStatement(cmd)))
}

func getStatement(cmd redis.Cmder) string {
	n := len(cmd.Args())
	if n > 1 {
		sb := &strings.Builder{}
		sb.WriteString(cmd.Name())
		sb.WriteString(" ")
		last := n - 2
		for i, arg := range cmd.Args()[1:] {
			sb.WriteString(fmt.Sprint(arg))
			if i < last {
				sb.WriteString(" ")
			}
		}
		return sb.String()
	}
	return cmd.Name()
}

var statementsCountKey = attribute.Key("db.statements_count")
var isBatchStatementsKey = attribute.Key("db.is_batch_statement")

func withBatchDBStatement(cfg *config, cmds []redis.Cmder) []attribute.KeyValue {
	if cfg.SpanOptions.DisableStatement {
		return cfg.Attributes
	}
	var statement string
	n := len(cmds)
	switch n {
	case 0:
	case 1:
		statement = getStatement(cmds[0])
	case 2:
		statement = getStatement(cmds[0]) + "; " + getStatement(cmds[1])
	default:
		statement = getStatement(cmds[0]) + "; ... ;" + getStatement(cmds[n-1])
	}
	return append(cfg.Attributes,
		semconv.DBStatementKey.String(statement),
		isBatchStatementsKey.Bool(true),
		statementsCountKey.Int64(int64(n)),
	)
}
