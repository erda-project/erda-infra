package protocol

import (
	"fmt"
	"github.com/autom8ter/dagger"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/sirupsen/logrus"
	"sync"
)

const NodeStatusRunning = "running"
const NodeStatusDone = "done"

type ComponentRender struct {
	dag   *dagger.Graph
	items []cptype.RendingItem
	lock  sync.Mutex

	doneNode map[string]string
}

func NewComponentRender() *ComponentRender {
	var componentRender = ComponentRender{}
	componentRender.dag = dagger.NewGraph()
	componentRender.doneNode = map[string]string{}
	return &componentRender
}

func (c *ComponentRender) addEdge(item cptype.RendingItem) error {
	if len(item.State) == 0 {
		if len(c.items) > 0 {
			_, err := c.dag.SetEdge(dagger.Path{
				XType: "depends_on",
				XID:   c.items[0].Name,
			}, dagger.Path{
				XID:   item.Name,
				XType: "depends_on",
			}, dagger.Node{
				Path: dagger.Path{
					XType: "depends_on",
				},
				Attributes: map[string]interface{}{
					"doneItemName": c.items[0].Name,
				},
			})
			if err != nil {
				logrus.Errorf("failed to set dag edge, component: %s", item.Name)
				return err
			}
		}
	} else {
		for _, state := range item.State {
			stateFrom, _, err := parseStateBound(state.Value)
			if err != nil {
				logrus.Errorf("failed to parse component state bound, component: %s, state bound: %#v", item.Name, state)
				return err
			}
			switch stateFrom {
			case cptype.InParamsStateBindingKey:
				continue
			default:
				_, err = c.dag.SetEdge(dagger.Path{
					XType: "depends_on",
					XID:   stateFrom,
				}, dagger.Path{
					XID:   item.Name,
					XType: "depends_on",
				}, dagger.Node{
					Path: dagger.Path{
						XType: "depends_on",
					},
					Attributes: map[string]interface{}{
						"doneItemName": stateFrom,
					},
				})
				if err != nil {
					logrus.Errorf("failed to set dag edge, component: %s, state bound: %#v", item.Name, state)
					return err
				}
			}
		}
	}
	return nil
}

func (c *ComponentRender) addNode(item cptype.RendingItem) error {
	c.dag.SetNode(dagger.Path{
		XID:   item.Name,
		XType: "depends_on",
	}, map[string]interface{}{
		"item": item,
	})

	c.items = append(c.items, item)
	return nil
}

func (c *ComponentRender) render(name string, doing func(name string) error) error {
	var dagEdges []dagger.Edge
	c.dag.RangeEdgesFrom("*", dagger.Path{
		XType: "depends_on",
		XID:   name,
	}, func(dagEdge dagger.Edge) bool {
		dagEdges = append(dagEdges, dagEdge)
		return true
	})

	var wait sync.WaitGroup
	var warn error
	for _, v := range dagEdges {
		if v.From.XID == v.To.XID {
			continue
		}
		wait.Add(1)
		go func(v dagger.Edge) {
			defer wait.Done()

			c.lock.Lock()
			var allDone = true
			c.dag.RangeEdgesTo("*", dagger.Path{
				XID:   v.To.XID,
				XType: "depends_on",
			}, func(e dagger.Edge) bool {
				done := c.doneNode[e.From.XID]
				if done != NodeStatusDone {
					allDone = false
				}
				return true
			})
			c.lock.Unlock()
			if !allDone {
				return
			}

			c.lock.Lock()
			done := c.doneNode[v.To.XID]
			if done != "" {
				c.lock.Unlock()
				return
			}
			c.doneNode[v.To.XID] = NodeStatusRunning
			c.lock.Unlock()

			err := doing(v.To.XID)
			if err != nil {
				warn = err
				return
			}

			c.lock.Lock()
			c.doneNode[v.To.XID] = NodeStatusDone
			c.lock.Unlock()

			fmt.Println(v.From.XID, "_", v.To.XID)

			err = c.render(v.To.XID, doing)
			if err != nil {
				warn = err
				return
			}
		}(v)
	}

	wait.Wait()
	return nil
}
