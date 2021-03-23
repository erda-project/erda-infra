// Author: recallsong
// Email: songruiguo@qq.com

package kafka

import (
	"reflect"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	writer "github.com/erda-project/erda-infra/pkg/parallel-writer"
)

// Interface .
type Interface interface {
	NewConsumer(c *ConsumerConfig, handler ConsumerFunc, options ...ConsumerOption) error
	NewConsumerWitchCreator(c *ConsumerConfig, creator func(i int) (ConsumerFunc, error), options ...ConsumerOption) error
	NewProducer(c *ProducerConfig, options ...ProducerOption) (writer.Writer, error)
	Servers() string
	ProduceChannelSize() int
	ProduceEventsChannelSize() int
	NewAdminClient() (*kafka.AdminClient, error)
}

// Producer .
type Producer interface {
	writer.Writer
	ProduceChannelSize() int
	EventsChannelSize() int
}

type config struct {
	Servers  string `file:"servers" env:"BOOTSTRAP_SERVERS" default:"localhost:9092" desc:"kafka servers"`
	Producer struct {
		Options map[string]interface{} `file:"options"`
	} `file:"producer"`
	Comsumer struct {
		Options map[string]interface{} `file:"options"`
	} `file:"comsumer"`
	Admin struct {
		Options map[string]interface{} `file:"options"`
	} `file:"admin"`
}

type define struct{}

func (d *define) Services() []string {
	return []string{"kafka", "kafka-producer", "kafka-consumer"}
}
func (d *define) Types() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Interface)(nil)).Elem(),
	}
}
func (d *define) Description() string { return "kafka、kafka-producer、kafka-consumer" }
func (d *define) Config() interface{} { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

// provider .
type provider struct {
	Cfg      *config
	Log      logs.Logger
	producer sharedProducer
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.producer.log = p.Log
	return nil
}

// Provide .
func (p *provider) Provide(name string, args ...interface{}) interface{} {
	return &service{
		p:    p,
		log:  p.Log.Sub(name),
		name: name,
	}
}

type service struct {
	p    *provider
	log  logs.Logger
	name string
}

func (s *service) Servers() string { return s.p.Cfg.Servers }

func (s *service) NewAdminClient() (*kafka.AdminClient, error) {
	kc := kafka.ConfigMap{}
	if s.p.Cfg.Admin.Options != nil {
		for k, v := range s.p.Cfg.Admin.Options {
			kc[k] = v
		}
	}
	kc["bootstrap.servers"] = s.p.Cfg.Servers
	return kafka.NewAdminClient(&kc)
}

func init() {
	servicehub.RegisterProvider("kafka", &define{})
}
