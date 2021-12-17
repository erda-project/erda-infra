package top

import "github.com/erda-project/erda-infra/providers/component-protocol/cptype"

type ITop interface {
	cptype.IComponent
	ITopStdOps
}

type ITopStdOps interface{}
