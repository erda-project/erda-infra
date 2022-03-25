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

package clickhouse

import (
	"context"
	"fmt"

	ckdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// EncodeFunc .
type EncodeFunc func(data interface{}) (item *WriteItem, err error)

// WriteItem .
type WriteItem struct {
	Table string
	Data  interface{}
}

// WriterOptions .
type WriterOptions struct {
	Encoder EncodeFunc
}

// Writer .
type Writer struct {
	client  ckdriver.Conn
	Encoder EncodeFunc
}

// NewWriter .
func NewWriter(client ckdriver.Conn, encoder EncodeFunc) *Writer {
	w := &Writer{
		client:  client,
		Encoder: encoder,
	}
	return w
}

// Close .
func (w *Writer) Close() error {
	return nil
}

// WriteN .
func (w *Writer) WriteN(list ...interface{}) (int, error) {
	if len(list) <= 0 {
		return 0, nil
	}

	items := map[string][]*WriteItem{}
	for _, data := range list {
		item, err := w.Encoder(data)
		if err != nil {
			return 0, err
		}

		items[item.Table] = append(items[item.Table], item)
	}

	succ := 0
	for table, tItems := range items {
		batch, err := w.client.PrepareBatch(context.Background(), fmt.Sprintf("insert into %s", table))
		if err != nil {
			return succ, err
		}
		for _, item := range tItems {
			err = batch.AppendStruct(item.Data)
			if err != nil {
				_ = batch.Abort()
				return succ, err
			}
		}
		err = batch.Send()
		if err != nil {
			return succ, err
		}
		succ++
	}

	return succ, nil
}
