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

package sqlite3

type Options struct {
	JournalMode JournalMode
	RandomName  bool
}

type OptionFunc func(options *Options)

type JournalMode string

const (
	DELETE   JournalMode = "delete"
	TRUNCATE JournalMode = "truncate"
	PERSIST  JournalMode = "persist"
	MEMORY   JournalMode = "memory"
	OFF      JournalMode = "off"
	WAL      JournalMode = "wal"
)

func WithJournalMode(mode JournalMode) OptionFunc {
	return func(o *Options) {
		o.JournalMode = mode
	}
}

// WithRandomName use to set the uuid in the given filename
// such as `test.db => test-550e8400e29b41d4a716446655440000.db`
func WithRandomName(isRandomName bool) OptionFunc {
	return func(o *Options) {
		o.RandomName = isRandomName
	}
}
