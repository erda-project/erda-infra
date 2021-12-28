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

package kubernetes

import (
	"context"
	"crypto/sha1" // nolint
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/kubernetes/watcher"
	"github.com/erda-project/erda-infra/providers/kubernetes/watcher/pod"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	certutil "k8s.io/client-go/util/cert"
)

// Interface .
type Interface interface {
	Client() *kubernetes.Clientset
	WatchPod(ctx context.Context, log logs.Logger, selector watcher.Selector) <-chan pod.Event
}

var clientType = reflect.TypeOf((*kubernetes.Clientset)(nil))

type config struct {
	MasterURL  string `file:"master_url"`
	ConfigPath string `file:"config_path"`
	RootCAFile string `file:"root_ca_file"`
	TokenFile  string `file:"token_file"`
}

// provider .
type provider struct {
	Cfg    *config
	Log    logs.Logger
	client *kubernetes.Clientset
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	if len(p.Cfg.ConfigPath) <= 0 {
		if home := homeDir(); home != "" {
			p.Cfg.ConfigPath = filepath.Join(home, ".kube", "config")
		}
	}
	config, err := p.createRestConfig()
	if err != nil {
		return fmt.Errorf("create rest config err: %w", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create k8s client err: %w", err)
	}

	if err := HealthCheck(clientset, 30*time.Second); err != nil {
		return fmt.Errorf("check connection err: %w", err)
	}

	p.client = clientset
	return nil
}

func (p *provider) WatchPod(ctx context.Context, log logs.Logger, selector watcher.Selector) <-chan pod.Event {
	w := pod.NewWatcher(
		ctx,
		p.client,
		p.Log.Sub(fmt.Sprintf("pod-"+identifySelector(selector))),
		watcher.Selector{
			Namespace:     selector.Namespace,
			LabelSelector: selector.LabelSelector,
			FieldSelector: selector.FieldSelector,
		},
	)
	ch := make(chan pod.Event)
	go w.Watch(ctx, ch)
	return ch
}

func (p *provider) Client() *kubernetes.Clientset { return p.client }

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Type() == clientType || ctx.Service() == "kube-client" || ctx.Service() == "kubernetes-client" {
		return p.client
	}
	return p
}

func (p *provider) createRestConfig() (*rest.Config, error) {
	var config *rest.Config
	if p.Cfg.MasterURL != "" {
		if p.Cfg.ConfigPath != "" {
			if _, err := os.Stat(p.Cfg.ConfigPath); err == nil {
				// use the current context in kubeconfig
				cfg, err := clientcmd.BuildConfigFromFlags(p.Cfg.MasterURL, p.Cfg.ConfigPath)
				if err != nil {
					return nil, fmt.Errorf("fail to build kube config: %s", err)
				}
				config = cfg
			}
		} else if p.Cfg.RootCAFile != "" && p.Cfg.TokenFile != "" {
			tlscfg := rest.TLSClientConfig{}
			if _, err := certutil.NewPool(p.Cfg.RootCAFile); err != nil {
				return nil, fmt.Errorf("expected to load root CA config from %s, but got err: %v", p.Cfg.RootCAFile, err)
			} else {
				tlscfg.CAFile = p.Cfg.RootCAFile
			}

			token, err := ioutil.ReadFile(p.Cfg.TokenFile)
			if err != nil {
				return nil, err
			}

			config = &rest.Config{
				TLSClientConfig: tlscfg,
				Host:            p.Cfg.MasterURL,
				BearerTokenFile: p.Cfg.TokenFile,
				BearerToken:     string(token),
			}
		}
	} else {
		cfg, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("build from inCluster err: %w", err)
		}
		config = cfg
	}
	return config, nil
}

func HealthCheck(client *kubernetes.Clientset, to time.Duration) error {
	ctx, _ := context.WithTimeout(context.TODO(), to)
	_, err := client.Discovery().RESTClient().Get().AbsPath("/healthz").DoRaw(ctx)
	return err
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func identifySelector(selector watcher.Selector) string {
	h := sha1.New()
	h.Write([]byte(strings.Join([]string{selector.Namespace, selector.LabelSelector, selector.FieldSelector}, ":")))
	return string(h.Sum(nil))[:10]
}

func init() {
	servicehub.Register("kubernetes", &servicehub.Spec{
		Services: []string{"kubernetes", "kubernetes-client", "kube-client"},
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
			clientType,
		},
		Description: "kubernetes",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
