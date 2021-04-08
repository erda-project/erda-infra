// Copyright (c) 2021 Terminus, Inc.

// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later (AGPL), as published by the Free Software Foundation.

// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.

// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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
