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
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/gocql/gocql"
)

type StatementBuilder interface {
	GetStatement(data interface{}) (string, []interface{}, error)
}

type batchWriter struct {
	session       *gocql.Session
	builder       StatementBuilder
	retry         int
	retryDuration time.Duration
	log           logs.Logger
}

func (w *batchWriter) Write(data interface{}) error {
	_, err := w.WriteN(data)
	return err
}

func (w *batchWriter) WriteN(data ...interface{}) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	batch := w.session.NewBatch(gocql.LoggedBatch)
	for _, item := range data {
		stmt, args, err := w.builder.GetStatement(item)
		if err != nil {
			w.log.Errorf("fail to convert data to statement: %s", err)
			continue
		}
		batch.Query(stmt, args...)
	}
	for i := 0; ; i++ {
		err := w.session.ExecuteBatch(batch)
		if err != nil {
			if w.retry == -1 || i < w.retry {
				w.log.Warnf("fail to write batch(%d) to cassandra and retry after %s: %s", batch.Size(), w.retryDuration.String(), err)
				time.Sleep(w.retryDuration)
				continue
			}
			w.log.Errorf("fail to write batch(%d) to cassandra: %s", batch.Size(), err)
			break
		}
		break
	}
	return batch.Size(), nil
}

func (w *batchWriter) Close() error {
	if w.session == nil {
		return nil
	}
	w.session.Close()
	w.session = nil
	return nil
}
