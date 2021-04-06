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

package servicehub

import (
	"bytes"
	"reflect"
)

// Usage .
func Usage(names ...string) string {
	buf := &bytes.Buffer{}
	buf.WriteString("Service Providers:\n")
	if len(names) <= 0 {
		for name, define := range serviceProviders {
			providerUsage(name, define, buf)
		}
	} else {
		for _, name := range names {
			if define, ok := serviceProviders[name]; ok {
				providerUsage(name, define, buf)
			}
		}
	}
	return buf.String()
}

func providerUsage(name string, define ProviderDefine, buf *bytes.Buffer) {
	buf.WriteString(name)
	var usage string
	if s, ok := define.(ProviderUsageSummary); ok {
		usage = s.Summary()
	}
	if len(usage) <= 0 {
		if u, ok := define.(ProviderUsage); ok {
			usage = u.Description()
		}
	}
	if len(usage) > 0 {
		buf.WriteString("\n    ")
		buf.WriteString(usage)
	}
	if creator, ok := define.(ConfigCreator); ok {
		cfg := creator.Config()
		typ := reflect.TypeOf(cfg)
		for typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		num := typ.NumField()
		for i := 0; i < num; i++ {
			field := typ.Field(i)
			file := field.Tag.Get("file")
			flag := field.Tag.Get("flag")
			env := field.Tag.Get("env")
			defval := field.Tag.Get("default")
			desc := field.Tag.Get("desc")
			var tags []string
			if len(file) > 0 {
				tags = append(tags, "file:\""+file+"\"")
			}
			if len(flag) > 0 {
				tags = append(tags, "flag:\""+flag+"\"")
			}
			if len(env) > 0 {
				tags = append(tags, "env:\""+env+"\"")
			}
			if len(defval) > 0 {
				tags = append(tags, "default:\""+defval+"\"")
			}
			if len(desc) > 0 {
				tags = append(tags, ", "+desc)
			}
			buf.WriteString("\n    ")
			for _, tag := range tags {
				buf.WriteString(tag)
				buf.WriteRune(' ')
			}
		}
	}
	buf.WriteRune('\n')
}
