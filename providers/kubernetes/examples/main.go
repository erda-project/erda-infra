// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"
	pkube "github.com/erda-project/erda-infra/providers/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type define struct{}

func (d *define) Service() []string      { return []string{"hello"} }
func (d *define) Dependencies() []string { return []string{"kubernetes"} }
func (d *define) Description() string    { return "hello for example" }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
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
	return nil
}

func init() {
	servicehub.RegisterProvider("examples", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
