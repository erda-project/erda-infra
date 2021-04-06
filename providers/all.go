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

package providers

import (
	_ "github.com/erda-project/erda-infra/providers/elasticsearch"        //
	_ "github.com/erda-project/erda-infra/providers/etcd"                 //
	_ "github.com/erda-project/erda-infra/providers/etcd-mutex"           //
	_ "github.com/erda-project/erda-infra/providers/grpcclient"           //
	_ "github.com/erda-project/erda-infra/providers/grpcserver"           //
	_ "github.com/erda-project/erda-infra/providers/health"               //
	_ "github.com/erda-project/erda-infra/providers/httpserver"           //
	_ "github.com/erda-project/erda-infra/providers/i18n"                 //
	_ "github.com/erda-project/erda-infra/providers/kubernetes"           //
	_ "github.com/erda-project/erda-infra/providers/legacy/httpendpoints" //
	_ "github.com/erda-project/erda-infra/providers/mysql"                //
	_ "github.com/erda-project/erda-infra/providers/pprof"                //
	_ "github.com/erda-project/erda-infra/providers/redis"                //
	_ "github.com/erda-project/erda-infra/providers/serviceregister"      //
	// _ "github.com/erda-project/erda-infra/providers/zk-master-election"   //
	// _ "github.com/erda-project/erda-infra/providers/zookeeper"            //
)
