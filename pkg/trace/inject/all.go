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

package inject

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	_ "github.com/erda-project/erda-infra/pkg/trace/inject/etcd-clientv3" //nolint
	_ "github.com/erda-project/erda-infra/pkg/trace/inject/http-client"   //nolint
	_ "github.com/erda-project/erda-infra/pkg/trace/inject/http-server"   //nolint
	_ "github.com/erda-project/erda-infra/pkg/trace/inject/redis"         //nolint
	_ "github.com/erda-project/erda-infra/pkg/trace/inject/sql"           //nolint
)

// Init .
func Init(opt ...sdktrace.TracerProviderOption) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERROR] failed to initialize TracerProvider %s\n", err)
		}
	}()
	tp := sdktrace.NewTracerProvider(opt...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}
