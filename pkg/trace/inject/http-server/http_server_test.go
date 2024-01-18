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

package traceinject

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHookHttpServer(t *testing.T) {
	// start a http server listen on 8083
	http.HandleFunc("/test",
		func(w http.ResponseWriter, r *http.Request) {
			h := getServerHandler(r.Context()) // injected by hook
			assert.NotNil(t, h)
			w.Write([]byte("hello world"))
		},
	)
	s := &http.Server{Addr: ":8083"}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(time.Second)
	// call api
	resp, err := http.Get("http://localhost:8083/test")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal("http status code not match")
	}
	// print body
	b, _ := io.ReadAll(resp.Body)
	t.Log(string(b))
	assert.Equal(t, "hello world", string(b)) // check original response
}
