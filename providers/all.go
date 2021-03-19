package providers

import (
	_ "github.com/erda-project/erda-infra/providers/elasticsearch"        //
	_ "github.com/erda-project/erda-infra/providers/etcd"                 //
	_ "github.com/erda-project/erda-infra/providers/etcd-mutex"           //
	_ "github.com/erda-project/erda-infra/providers/health"               //
	_ "github.com/erda-project/erda-infra/providers/httpserver"           //
	_ "github.com/erda-project/erda-infra/providers/i18n"                 //
	_ "github.com/erda-project/erda-infra/providers/kubernetes"           //
	_ "github.com/erda-project/erda-infra/providers/legacy/httpendpoints" //
	_ "github.com/erda-project/erda-infra/providers/mysql"                //
	_ "github.com/erda-project/erda-infra/providers/pprof"                //
	_ "github.com/erda-project/erda-infra/providers/redis"                //
	_ "github.com/erda-project/erda-infra/providers/zk-master-election"   //
	_ "github.com/erda-project/erda-infra/providers/zookeeper"            //
)
