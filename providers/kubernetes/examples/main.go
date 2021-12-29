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
	"fmt"
	"os"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/providers/kubernetes/watcher"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/erda-project/erda-infra/base/servicehub"
	pkube "github.com/erda-project/erda-infra/providers/kubernetes"
)

type provider struct {
	Log    logs.Logger
	Kube   pkube.Interface
	Client *kubernetes.Clientset
}

func (p *provider) Init(ctx servicehub.Context) error {
	// fmt.Println(p.Kube)
	// fmt.Println(p.Client)
	return nil
}

func (p *provider) Run(ctx context.Context) error {
	nodes, err := p.Client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, node := range nodes.Items {
		ip := node.Name
		fmt.Println(ip)
	}

	ch := p.Kube.WatchPod(ctx, p.Log.Sub("pod-watch"), watcher.Selector{
		Namespace:     "default",
	})

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-ch:
			fmt.Println(event.Pod.Name)
		}
	}
}

func init() {
	servicehub.Register("examples", &servicehub.Spec{
		Services:     []string{"hello"},
		Dependencies: []string{"kubernetes"},
		Description:  "hello for example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
