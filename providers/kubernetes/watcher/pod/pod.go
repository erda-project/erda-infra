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

package pod

import (
	"context"
	"fmt"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/providers/kubernetes/watcher"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Watcher struct {
	informer cache.SharedInformer
	store    cache.Store
	queue    *workqueue.Type
	log      logs.Logger
}

type Event struct {
	Pod    *apiv1.Pod
	Action watcher.Action
}

func NewWatcher(ctx context.Context, c *k8s.Clientset, log logs.Logger, selector watcher.Selector) *Watcher {
	pg := c.CoreV1().Pods(selector.Namespace)
	informer := cache.NewSharedInformer(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = selector.FieldSelector
			options.LabelSelector = selector.LabelSelector
			return pg.List(ctx, options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = selector.FieldSelector
			options.LabelSelector = selector.LabelSelector
			return pg.Watch(ctx, options)
		},
	}, &apiv1.Pod{}, 10*time.Minute)
	go informer.Run(ctx.Done())

	p := &Watcher{
		informer: informer,
		store:    informer.GetStore(),
		queue:    workqueue.NewNamed("pod"),
		log:      log,
	}

	p.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			p.enqueue(obj, watcher.ActionAdd)
		},
		UpdateFunc: func(_, obj interface{}) {
			p.enqueue(obj, watcher.ActionAdd)
		},
		DeleteFunc: func(obj interface{}) {
			p.enqueue(obj, watcher.ActionDelete)
		},
	})
	return p
}

func (p *Watcher) enqueue(obj interface{}, action watcher.Action) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		return
	}
	p.queue.Add(&watcher.Item{Key: key, Action: action, Object: obj})
}

func (p *Watcher) Watch(ctx context.Context, ch chan<- Event) {
	defer p.queue.ShutDown()
	if !cache.WaitForCacheSync(ctx.Done(), p.informer.HasSynced) {
		if ctx.Err() != context.Canceled {
			p.log.Info("pod informer unable to sync cache")
		}
		return
	}

	go func() {
		for {
			ok, err := p.process(ctx, ch)
			if err != nil {
				p.log.Errorf("process err: %s", err)
			}
			if !ok {
				break
			}
		}
	}()

	<-ctx.Done()
}

func (p *Watcher) process(ctx context.Context, ch chan<- Event) (bool, error) {
	keyObj, quit := p.queue.Get()
	if quit {
		return false, nil
	}
	defer p.queue.Done(keyObj)

	item, ok := keyObj.(*watcher.Item)
	if !ok {
		return true, fmt.Errorf("convert to *watcher.Item error")
	}

	_, exist, err := p.store.GetByKey(item.Key)
	if err != nil {
		return true, fmt.Errorf("get key %q, err: %w", item.Key, err)
	}
	if !exist && item.Action != watcher.ActionDelete {
		return true, fmt.Errorf("key not existed in store. item: %+v", item)
	}

	pod, ok := item.Object.(*apiv1.Pod)
	if !ok {
		return true, fmt.Errorf("convert to *apiv1.Pod error")
	}
	e := Event{
		Pod:    pod,
		Action: item.Action,
	}
	select {
	case <-ctx.Done():
		return false, nil
	case ch <- e:
		return true, nil
	}
}
