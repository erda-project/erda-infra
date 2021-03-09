// Author: recallsong
// Email: songruiguo@qq.com

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
	if usage, ok := define.(ProviderUsageSummary); ok {
		buf.WriteString("\n    ")
		buf.WriteString(usage.Summary())
	} else if usage, ok := define.(ProviderUsage); ok {
		buf.WriteString("\n    ")
		buf.WriteString(usage.Description())
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
