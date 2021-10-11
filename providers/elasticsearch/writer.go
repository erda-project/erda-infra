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

package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/olivere/elastic"
	"github.com/recallsong/go-utils/reflectx"
)

// Document .
type Document struct {
	ID    string      `json:"id"`
	Index string      `json:"index"`
	Data  interface{} `json:"data"`
}

type batchWriter struct {
	client        *elastic.Client
	typ           string
	timeout       string
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

	requests := make([]elastic.BulkableRequest, len(data), len(data))
	for i, item := range data {
		if doc, ok := item.(*Document); ok {
			req := elastic.NewBulkIndexRequest().Index(doc.Index).Type(w.typ).Doc(doc.Data)
			if len(doc.ID) > 0 {
				req.Id(doc.ID)
			}
			requests[i] = req
		} else {
			return 0, fmt.Errorf("%s is not *elasticsearch.Document", reflect.TypeOf(item))
		}
	}
	for i := 0; ; i++ {
		res, err := w.client.Bulk().Add(requests...).Timeout(w.timeout).Do(context.Background())
		if err != nil {
			if i < w.retry {
				w.log.Warnf("failed to write batch(%d) to elasticsearch and retry after %s: %s", len(requests), w.retryDuration.String(), err)
				time.Sleep(w.retryDuration)
				continue
			}
			w.log.Errorf("failed to write batch(%d) to elasticsearch: %s", len(requests), err)
			break
		}
		if res.Errors {
			for _, item := range res.Failed() {
				if item == nil || item.Error == nil {
					continue
				}
				byts, _ := json.Marshal(item.Error)
				// TODO: notify error
				w.log.Errorf("failed to index data, [%s][%s]: %s", item.Index, item.Type, reflectx.BytesToString(byts))
			}
		}
		break
	}
	return len(requests), nil
}

func (w *batchWriter) Close() error { return nil }

// NewWriter .
func NewWriter(client *elastic.Client, timeout time.Duration, enc EncodeFunc) *Writer {
	w := &Writer{
		client: client,
		enc:    enc,
	}
	if timeout > 0 {
		w.timeout = fmt.Sprintf("%dms", timeout.Milliseconds())
	}
	return w
}

// Writer .
type Writer struct {
	client  *elastic.Client
	enc     EncodeFunc
	timeout string
}

func (w *Writer) Close() error { return nil }

func (w *Writer) Write(data interface{}) error {
	index, id, typ, body := w.enc(data)
	_, err := w.client.Index().
		Index(index).Id(id).Type(typ).
		BodyJson(body).Timeout(w.timeout).Do(context.Background())
	return err
}

func (w *Writer) WriteN(list ...interface{}) (int, error) {
	if len(list) <= 0 {
		return 0, nil
	}
	bulk := w.client.Bulk()
	for _, data := range list {
		index, id, typ, body := w.enc(data)
		req := elastic.NewBulkIndexRequest().Index(index).Id(id).Type(typ).Doc(body)
		bulk.Add(req)
	}
	res, err := bulk.Timeout(w.timeout).Do(context.Background())
	if err != nil {
		berr := &BatchWriteError{
			List:   list,
			Errors: make([]error, len(list), len(list)),
		}
		for i, n := 0, len(list); i < n; i++ {
			berr.Errors[i] = err
		}
		return 0, berr
	}
	if res.Errors {
		if len(res.Items) != len(list) {
			return 0, fmt.Errorf("request items(%d), but response items(%d)", len(list), len(res.Items))
		}
		berr := &BatchWriteError{
			List:   make([]interface{}, 0, len(list)),
			Errors: make([]error, 0, len(list)),
		}
		for i, item := range res.Items {
			for _, result := range item { // len(item) is 1, contains index request only
				if !(result.Status >= 200 && result.Status <= 299) {
					var sb strings.Builder
					json.NewEncoder(&sb).Encode(result)
					berr.List = append(berr.List, list[i])
					berr.Errors = append(berr.Errors, errors.New(sb.String()))
					break
				}
			}
		}
		return len(list) - len(berr.Errors), berr
	}
	return len(list), nil
}

// BatchWriteError .
type BatchWriteError struct {
	List   []interface{}
	Errors []error
}

func (e *BatchWriteError) Error() string {
	if len(e.Errors) <= 0 {
		return ""
	}
	return fmt.Sprintf("bulk writes occur errors(%d): %v ...", len(e.Errors), e.Errors[0])
}
