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

// StatementBuilder .
type StatementBuilder interface {
	GetStatement(data interface{}) (string, []interface{}, error)
}

type batchWriter struct {
	session        *Session
	builder        StatementBuilder
	retry          int
	retryDuration  time.Duration
	log            logs.Logger
	batchSizeBytes int
}

func (w *batchWriter) Write(data interface{}) error {
	_, err := w.WriteN(data)
	return err
}

func (w *batchWriter) WriteN(data ...interface{}) (int, error) {
	if len(data) <= 0 {
		return 0, nil
	}
	batchs := make([]*gocql.Batch, 0, 1)
	sizeBytes := 0

	batch := w.session.Session().NewBatch(gocql.LoggedBatch)
	for _, item := range data {
		stmt, args, err := w.builder.GetStatement(item)
		if err != nil {
			w.log.Errorf("fail to convert data to statement: %s", err)
			continue
		}
		if w.batchSizeBytes > 0 {
			sizeBytes += cqlSizeBytes(stmt, args)
			if sizeBytes >= w.batchSizeBytes {
				batchs = append(batchs, batch)
				// reset
				sizeBytes = 0
				batch = w.session.Session().NewBatch(gocql.LoggedBatch)
			}

		}
		batch.Query(stmt, args...)
	}

	if batch.Size() > 0 {
		batchs = append(batchs, batch)
	}

	for _, batch := range batchs {
		for i := 0; ; i++ {
			err := w.session.Session().ExecuteBatch(batch)
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
	}

	return batch.Size(), nil
}

func cqlSizeBytes(stmt string, args []interface{}) int {
	size := len(stmt)
	for _, item := range args {
		switch v := item.(type) {
		case string:
			size += len(v)
		case []byte:
			size += len(v)
		}
	}
	return size
}

func (w *batchWriter) Close() error {
	if w.session == nil {
		return nil
	}
	w.session.Close()
	w.session = nil
	return nil
}
