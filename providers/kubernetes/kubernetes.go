// Author: recallsong
// Email: songruiguo@qq.com

package kubernetes

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Interface .
type Interface interface {
	Client() *kubernetes.Clientset
}

var clientType = reflect.TypeOf((*kubernetes.Clientset)(nil))

type define struct{}

func (d *define) Services() []string {
	return []string{"kubernetes", "kubernetes-client", "kube-client"}
}
func (d *define) Types() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*Interface)(nil)).Elem(),
		clientType,
	}
}
func (d *define) Description() string { return "kubernetes" }
func (d *define) Config() interface{} { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type config struct {
	ConfigPath string `file:"config_path"`
	MasterURL  string `file:"master_url"`
}

// provider .
type provider struct {
	Cfg    *config
	client *kubernetes.Clientset
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	if len(p.Cfg.ConfigPath) <= 0 {
		if home := homeDir(); home != "" {
			p.Cfg.ConfigPath = filepath.Join(home, ".kube", "config")
		}
	}
	if len(p.Cfg.ConfigPath) <= 0 && len(p.Cfg.MasterURL) <= 0 {
		return fmt.Errorf("kube config path or master url must not be empty")
	}
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags(p.Cfg.MasterURL, p.Cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("fail to build kube config: %s", err)
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("fail to create k8s client: %s", err)
	}
	p.client = clientset
	return nil
}

func (p *provider) Client() *kubernetes.Clientset { return p.client }

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Type() == clientType || ctx.Service() == "kube-client" || ctx.Service() == "kubernetes-client" {
		return p.client
	}
	return p
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func init() {
	servicehub.RegisterProvider("kubernetes", &define{})
}
