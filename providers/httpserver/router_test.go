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

package httpserver

import (
	"testing"
	"time"

	"github.com/go-playground/validator"

	"github.com/erda-project/erda-infra/base/logs/logrusx"
	"github.com/erda-project/erda-infra/providers/httpserver/server"
)

func Test_removeParamName(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			path: "/abc",
			want: "/abc",
		},
		{
			path: "/abc/:def",
			want: "/abc/*",
		},
		{
			path: "/abc/:def/ghi",
			want: "/abc/*/ghi",
		},
		{
			path: "/ab/:c/*/:de/f/**",
			want: "/ab/*/*/*/f/**",
		},
		{
			path: "/ab/:c/*/:de/f/**/:G",
			want: "/ab/*/*/*/f/**/*",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeParamName(tt.path); got != tt.want {
				t.Errorf("removeParamName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_routerManager_Started(t *testing.T) {
	srv := server.New(false, &dataBinder{}, &structValidator{validator: validator.New()})
	p := &provider{server: srv, Cfg: &config{PrintRoutes: false}, Log: logrusx.New(), startedChan: make(chan struct{})}
	rm := &routerManager{p: p}
	go func() {
		// blocking http server start
		_ = p.Start()
	}()
	time.Sleep(time.Second)
	select {
	case <-rm.Started():
		// started
	default:
		t.Fatalf("should started")
	}
}
