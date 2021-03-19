module github.com/erda-project/erda-infra

go 1.14

require (
	github.com/confluentinc/confluent-kafka-go v1.6.1
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/go-zookeeper/zk v1.0.2
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/olivere/elastic v6.2.35+incompatible
	github.com/onsi/ginkgo v1.15.2 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/recallsong/go-utils v1.1.1
	github.com/recallsong/unmarshal v1.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/valyala/fasttemplate v1.2.1 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20210314154223-e6e6c4f2bb5b // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	google.golang.org/genproto v0.0.0-20210317182105-75c7a8546eb9 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.0.0-00010101000000-000000000000
)

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
	k8s.io/api => github.com/kubernetes/api v0.18.2
	k8s.io/apiextensions-apiserver => github.com/kubernetes/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery => github.com/kubernetes/apimachinery v0.18.2
	k8s.io/apiserver => github.com/kubernetes/apiserver v0.18.3
	k8s.io/client-go => github.com/kubernetes/client-go v0.18.2
	k8s.io/component-base => k8s.io/component-base v0.18.3
	k8s.io/klog => github.com/kubernetes/klog v1.0.0
	k8s.io/kube-scheduler => github.com/kubernetes/kube-scheduler v0.18.2
	k8s.io/kubectl => github.com/kubernetes/kubectl v0.18.3
	k8s.io/kubernetes => github.com/kubernetes/kubernetes v1.13.5
)
