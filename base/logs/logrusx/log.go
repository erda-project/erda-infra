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

package logrusx

import (
	"io"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/base/logs"
)

// Logger .
type Logger struct {
	name string
	*logrus.Entry
}

// New .
func New(options ...Option) logs.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	logger := &Logger{"", logrus.NewEntry(log)}
	for _, opt := range options {
		processOptions(log, logger, opt.get())
	}
	return logger
}

// Sub .
func (l *Logger) Sub(name string) logs.Logger {
	if len(l.name) > 0 {
		name = l.name + "." + name
	}
	return &Logger{name, l.Entry.WithField("module", name)}
}

// SetLevel .
func (l *Logger) SetLevel(lvl string) error {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return err
	}
	l.Logger.SetLevel(level)
	return nil
}

func (l *Logger) SetOutput(output io.Writer) {
	l.Logger.SetOutput(output)
}

func (l *Logger) AddHook(hook logrus.Hook) {
	l.Logger.AddHook(hook)
}

func (l *Logger) SetReportCaller(reportCaller bool) {
	l.Logger.SetReportCaller(reportCaller)
}

func (l *Logger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.SetFormatter(formatter)
}

func (l *Logger) SetNoLock() {
	l.Logger.SetNoLock()
}

func processOptions(logr *logrus.Logger, logger *Logger, opt interface{}) {
	switch val := opt.(type) {
	case setNameOption:
		logger.name = string(val)
	case logrus.Level:
		logr.SetLevel(val)
	}
}

// Option .
type Option interface {
	get() interface{}
}

type option struct{ value interface{} }

func (o *option) get() interface{} { return o.value }

type setNameOption string

// WithName .
func WithName(name string) Option {
	return &option{setNameOption(name)}
}

// WithLevel .
func WithLevel(level logrus.Level) Option {
	return &option{level}
}
