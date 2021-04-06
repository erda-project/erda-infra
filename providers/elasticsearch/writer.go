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

package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
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
	requests := make([]elastic.BulkableRequest, 0, len(data))
	for _, item := range data {
		if doc, ok := item.(*Document); ok {
			req := elastic.NewBulkIndexRequest().Index(doc.Index).Type(w.typ).Doc(doc.Data)
			if len(doc.ID) > 0 {
				req.Id(doc.ID)
			}
			requests = append(requests, req)
		} else {
			return 0, fmt.Errorf("%s is not *elasticsearch.Document", reflect.TypeOf(item))
		}
	}
	for i := 0; ; i++ {
		res, err := w.client.Bulk().Add(requests...).Timeout(w.timeout).Do(context.Background())
		if err != nil {
			if i < w.retry {
				w.log.Warnf("fail to write batch(%d) to elasticsearch and retry after %s: %s", len(requests), w.retryDuration.String(), err)
				time.Sleep(w.retryDuration)
				continue
			}
			w.log.Errorf("fail to write batch(%d) to elasticsearch: %s", len(requests), err)
			break
		}
		if res.Errors {
			for _, item := range res.Failed() {
				if item == nil || item.Error == nil {
					continue
				}
				byts, _ := json.Marshal(item.Error)
				w.log.Errorf("fail to index data, [%s][%s]: %s", item.Index, item.Type, reflectx.BytesToString(byts))
			}
		}
		break
	}
	return len(requests), nil
}

func (w *batchWriter) Close() error { return nil }
