module github.com/erda-project/erda-infra

go 1.16

require (
	github.com/Chronokeeper/anyxml v0.0.0-20160530174208-54457d8e98c6 // indirect
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/CloudyKit/jet v2.1.2+incompatible // indirect
	github.com/agrison/go-tablib v0.0.0-20160310143025-4930582c22ee // indirect
	github.com/agrison/mxj v0.0.0-20160310142625-1269f8afb3b4 // indirect
	github.com/bndr/gotabulate v1.1.2 // indirect
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/confluentinc/confluent-kafka-go v1.5.2
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/erda-project/erda-proto-go v0.0.0-20210820073528-bcd612c2efb7
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-xorm/sqlfiddle v0.0.0-20180821085327-62ce714f951a // indirect
	github.com/go-zookeeper/zk v1.0.2
	github.com/gocql/gocql v0.0.0-20210707082121-9a3953d1826d
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/labstack/echo v3.3.10+incompatible
	github.com/olivere/elastic v6.2.35+incompatible
	github.com/recallsong/go-utils v1.1.2-0.20210826100715-fce05eefa294
	github.com/recallsong/unmarshal v1.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tealeg/xlsx v1.0.5 // indirect
	github.com/xormplus/builder v0.0.0-20181220055446-b12ceebee76f // indirect
	github.com/xormplus/core v0.0.0-20181016121923-6bfce2eb8867
	github.com/xormplus/xorm v0.0.0-20181212020813-da46657160ff
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20210819135213-f52c844e1c1c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210820002220-43fce44e7af1
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/flosch/pongo2.v3 v3.0.0-20141028000813-5e81b817a0c4 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
)

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
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
