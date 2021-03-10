// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"os"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/i18n"
)

type define struct{}

func (d *define) Service() []string      { return []string{"hello"} }
func (d *define) Dependencies() []string { return []string{"i18n"} }
func (d *define) Description() string    { return "hello for example" }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{}
	}
}

type provider struct {
	L logs.Logger
}

func (p *provider) Init(ctx servicehub.Context) error {
	langs, err := i18n.ParseLanguageCode("en,zh-CN;q=0.9,zh;q=0.8,en-US;q=0.7,en-GB;q=0.6")
	if err != nil {
		return err
	}
	i := ctx.Service("i18n").(i18n.I18n)
	text := i.Text("hello", langs, "name")
	p.L.Info(text)
	return nil
}

func init() {
	servicehub.RegisterProvider("hello", &define{})
}

func main() {
	hub := servicehub.New()
	hub.Run("examples", "", os.Args...)
}
