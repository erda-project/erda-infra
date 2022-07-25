module github.com/erda-project/erda-infra

go 1.17

require (
	bou.ke/monkey v1.0.2
	github.com/Chronokeeper/anyxml v0.0.0-20160530174208-54457d8e98c6 // indirect
	github.com/ClickHouse/clickhouse-go/v2 v2.0.12
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet v2.1.2+incompatible // indirect
	github.com/XSAM/otelsql v0.9.0
	github.com/agrison/go-tablib v0.0.0-20160310143025-4930582c22ee // indirect
	github.com/agrison/mxj v0.0.0-20160310142625-1269f8afb3b4 // indirect
	github.com/bndr/gotabulate v1.1.2 // indirect
	github.com/brahma-adshonor/gohook v1.1.9
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/confluentinc/confluent-kafka-go v1.5.2
	github.com/coreos/bbolt v1.3.5 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-eden/routine v0.0.1
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-xorm/sqlfiddle v0.0.0-20180821085327-62ce714f951a // indirect
	github.com/go-zookeeper/zk v1.0.2
	github.com/gocql/gocql v0.0.0-20210707082121-9a3953d1826d
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/hcl v1.0.0
	github.com/hashicorp/yamux v0.0.0-20210826001029-26ff87cf9493
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.10
	github.com/kr/text v0.2.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.5
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/olivere/elastic v6.2.35+incompatible
	github.com/pelletier/go-toml v1.9.4
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.10.0
	github.com/recallsong/go-utils v1.1.2-0.20210826100715-fce05eefa294
	github.com/recallsong/unmarshal v1.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/kafka-go v0.4.31
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tealeg/xlsx v1.0.5 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	github.com/xormplus/builder v0.0.0-20181220055446-b12ceebee76f // indirect
	github.com/xormplus/core v0.0.0-20181016121923-6bfce2eb8867
	github.com/xormplus/xorm v0.0.0-20181212020813-da46657160ff
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.etcd.io/bbolt v1.3.5 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.27.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.27.0
	go.opentelemetry.io/otel v1.4.1
	go.opentelemetry.io/otel/sdk v1.2.0
	go.opentelemetry.io/otel/trace v1.4.1
	golang.org/x/crypto v0.0.0-20210920023735-84f357641f63 // indirect
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210820002220-43fce44e7af1
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/flosch/pongo2.v3 v3.0.0-20141028000813-5e81b817a0c4 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/ini.v1 v1.63.2
	gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	gorm.io/driver/mysql v1.3.2
	gorm.io/driver/sqlite v1.3.1
	gorm.io/gorm v1.23.2
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
)

replace (
	github.com/coreos/bbolt v1.3.5 => go.etcd.io/bbolt v1.3.5
	go.etcd.io/bbolt v1.3.5 => github.com/coreos/bbolt v1.3.5
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.2.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/api => github.com/kubernetes/api v0.18.3
	k8s.io/apiextensions-apiserver => github.com/kubernetes/apiextensions-apiserver v0.18.3
	k8s.io/apimachinery => github.com/kubernetes/apimachinery v0.18.3
	k8s.io/apiserver => github.com/kubernetes/apiserver v0.18.3
	k8s.io/client-go => github.com/kubernetes/client-go v0.18.3
	k8s.io/component-base => github.com/kubernetes/component-base v0.18.3
	k8s.io/klog => github.com/kubernetes/klog v1.0.0
	k8s.io/kube-scheduler => github.com/kubernetes/kube-scheduler v0.18.3
	k8s.io/kubectl => github.com/kubernetes/kubectl v0.18.3
	k8s.io/kubernetes => github.com/kubernetes/kubernetes v1.13.5
)
