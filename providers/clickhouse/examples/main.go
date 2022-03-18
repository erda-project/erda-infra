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

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	ck "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/clickhouse"
)

type provider struct {
	Clickhouse clickhouse.Interface
}

func (p *provider) Init(ctx servicehub.Context) error {
	// create table
	err := p.Clickhouse.Client().Exec(context.Background(),
		`create table if not exists example (
					timestamp DateTime(9,'Asia/Shanghai'),
					value String
				) Engine = Memory`)

	if err != nil {
		fmt.Println(err)
		return err
	}

	// batch insert
	batch, err := p.Clickhouse.Client().
		PrepareBatch(context.Background(), "insert into example")
	if err != nil {
		return err
	}
	err = batch.Append(time.Now(), "hello clickhouse")
	if err != nil {
		return err
	}
	err = batch.Send()
	if err != nil {
		return err
	}

	// query records
	// struct fields must be exposed
	// result fields binding is case-sensitive
	var result []struct {
		Timestamp time.Time `ch:"timestamp"`
		Value     string    `ch:"value"`
	}
	err = p.Clickhouse.Client().
		Select(context.Background(),
			&result,
			"select timestamp, value from example where timestamp < @time",
			ck.Named("time", time.Now()),
		)

	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("query result: %+v", result)

	// see https://github.com/ClickHouse/clickhouse-go/tree/v2/examples/native for more examples.
	return nil
}

func init() {
	servicehub.Register("example", &servicehub.Spec{
		Services:     []string{"example"},
		Dependencies: []string{"clickhouse"},
		Description:  "example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
