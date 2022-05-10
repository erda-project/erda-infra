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
	"embed"
	"os"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/i18n"
)

//go:embed i18n
var i18nFS embed.FS

type provider struct {
	Log  logs.Logger
	I18n i18n.I18n       `autowired:"i18n"`
	Tran i18n.Translator `translator:"hello"`
}

func (p *provider) Init(ctx servicehub.Context) error {
	if err := p.I18n.RegisterFilesFromFS("i18n", i18nFS); err != nil {
		return err
	}
	langs, err := i18n.ParseLanguageCode("en,zh-CN;q=0.9,zh;q=0.8,en-US;q=0.7,en-GB;q=0.6")
	if err != nil {
		return err
	}
	i := ctx.Service("i18n").(i18n.I18n)
	text := i.Text("hello", langs, "name")
	p.Log.Info(text)

	text = p.Tran.Text(langs, "name")
	p.Log.Info(text)

	text = p.Tran.Text(langs, "other")
	p.Log.Info(text)

	text = p.Tran.Text(langs, "common name")
	p.Log.Info(text)

	text = p.Tran.Text(langs, "file name")
	p.Log.Info(text)

	text = p.Tran.Sprintf(langs, "${Internal Error}: reason ${Reason} %s", "test error")
	p.Log.Info(text)

	text = i.Sprintf("not-exist-namespace", langs, "${Internal Error}: reason ${Reason} %s", "test error")
	p.Log.Info(text)
	return nil
}

func init() {
	servicehub.Register("hello", &servicehub.Spec{
		Services:     []string{"hello"},
		Dependencies: []string{"i18n"},
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

// OUTPUT:
// INFO[2021-09-14 18:05:25.629] load i18n files: [], [hello.yaml]             module=i18n
// INFO[2021-09-14 18:05:25.629] provider i18n initialized
// INFO[2021-09-14 18:05:25.629] 名字                                            module=hello
// INFO[2021-09-14 18:05:25.630] 名字                                            module=hello
// INFO[2021-09-14 18:05:25.630] other                                         module=hello
// INFO[2021-09-14 18:05:25.630] 内部错误: reason 未知 test error                    module=hello
// INFO[2021-09-14 18:05:25.630] Internal Error: reason Reason test error      module=hello
// INFO[2021-09-14 18:05:25.630] provider hello (depends services: [i18n], providers: [i18n]) initialized
// INFO[2021-09-14 18:05:25.630] signals to quit: [hangup interrupt terminated quit]
