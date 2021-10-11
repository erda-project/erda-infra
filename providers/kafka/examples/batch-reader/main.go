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
	"encoding/json"
	"fmt"
	"os"
	"time"

	gokafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/kafka"
)

type config struct {
	Input kafka.BatchReaderConfig `file:"input"`
}

type provider struct {
	Cfg   *config
	Log   logs.Logger
	Kafka kafka.Interface `autowired:"kafka-consumer"`
}

func (p *provider) Run(ctx context.Context) error {
	reader, err := p.Kafka.NewBatchReader(&p.Cfg.Input,
		kafka.WithReaderDecoder(func(key, value []byte, topic *string, timestamp time.Time) (interface{}, error) {
			m := make(map[string]interface{})
			err := json.Unmarshal(value, &m)
			return m, err
		}),
	)
	if err != nil {
		return err
	}
	defer reader.Close()

	bufferSize := 100
	limit := 200
	buf := make([]interface{}, bufferSize)

	for {
		// check
		if limit <= 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// read data
		n, err := reader.ReadN(buf, time.Second)
		if err != nil {
			return err
		}
		if n <= 0 {
			continue
		}

		// process data
		for i := 0; i < n; i++ {
			fmt.Println(buf[i])
			limit--
			if limit <= 0 {
				break
			}
		}

		// print offsets before confirm
		fmt.Println("offsets before confirm")
		ps, _ := kafka.CommittedOffsets(reader)
		printPartitions(ps)

		// commit offsets
		err = reader.Confirm()
		if err != nil {
			return nil
		}

		// print offsets after confirm
		fmt.Println("offsets after confirm")
		ps, _ = kafka.CommittedOffsets(reader)
		printPartitions(ps)
	}
}

func printPartitions(partitions []gokafka.TopicPartition) {
	for _, p := range partitions {
		byts, _ := json.Marshal(p)
		fmt.Println("partition:", string(byts))
	}
}

func init() {
	servicehub.Register("examples", &servicehub.Spec{
		Services:   []string{"hello"},
		ConfigFunc: func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
