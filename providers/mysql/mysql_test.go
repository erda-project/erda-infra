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

package mysql

import (
	"testing"
	"time"
)

func Test_config_url(t *testing.T) {
	type fields struct {
		MySQLURL          string
		MySQLHost         string
		MySQLPort         string
		MySQLUsername     string
		MySQLPassword     string
		MySQLDatabase     string
		MySQLMaxIdleConns uint64
		MySQLMaxOpenConns uint64
		MySQLMaxLifeTime  time.Duration
		MySQLDebug        bool
		MySQLCharset      string
		MySQLTLS          string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test_not_have_tls",
			fields: fields{
				MySQLUsername: "erda",
				MySQLHost:     "erda",
				MySQLCharset:  "erda",
				MySQLDatabase: "erda",
				MySQLPort:     "3306",
				MySQLPassword: "erda",
			},
			want: "erda:erda@tcp(erda:3306)/erda?charset=erda&parseTime=True&loc=Local",
		},
		{
			name: "test_not_have_tls",
			fields: fields{
				MySQLUsername: "erda",
				MySQLHost:     "erda",
				MySQLCharset:  "erda",
				MySQLDatabase: "erda",
				MySQLPort:     "3306",
				MySQLPassword: "erda",
				MySQLTLS:      "custom",
			},
			want: "erda:erda@tcp(erda:3306)/erda?charset=erda&parseTime=True&loc=Local&tls=custom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &config{
				MySQLURL:          tt.fields.MySQLURL,
				MySQLHost:         tt.fields.MySQLHost,
				MySQLPort:         tt.fields.MySQLPort,
				MySQLUsername:     tt.fields.MySQLUsername,
				MySQLPassword:     tt.fields.MySQLPassword,
				MySQLDatabase:     tt.fields.MySQLDatabase,
				MySQLMaxIdleConns: tt.fields.MySQLMaxIdleConns,
				MySQLMaxOpenConns: tt.fields.MySQLMaxOpenConns,
				MySQLMaxLifeTime:  tt.fields.MySQLMaxLifeTime,
				MySQLDebug:        tt.fields.MySQLDebug,
				MySQLCharset:      tt.fields.MySQLCharset,
				MySQLTLS:          tt.fields.MySQLTLS,
			}
			if got := c.url(); got != tt.want {
				t.Errorf("url() = %v, want %v", got, tt.want)
			}
		})
	}
}
