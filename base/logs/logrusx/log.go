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



package logrusx

import (
	"github.com/erda-project/erda-infra/base/logs"
	"github.com/sirupsen/logrus"
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
		ForceColors:     false,
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
