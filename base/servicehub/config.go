// Author: recallsong
// Email: songruiguo@qq.com

package servicehub

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/recallsong/go-utils/config"
)

func (h *Hub) loadConfigWithArgs(file string, args ...string) (map[string]interface{}, error) {
	cfg := make(map[string]interface{})
	if len(args) > 0 {
		args = args[1:]
		var idx int
		var arg string
		for idx, arg = range args {
			if !strings.HasPrefix(arg, "-") {
				cfg[arg] = nil
			} else {
				break
			}
		}
		args, num, idx := args[idx:], len(args[idx:]), 0
		for ; idx < num; idx++ {
			arg := args[idx]
			if strings.HasPrefix(arg, "-") {
				if arg == "-c" || arg == "--config" || arg == "-config" {
					idx++
					if idx < len(args) && len(args[idx]) > 0 {
						file = args[idx]
					}
				} else if strings.HasPrefix(arg, "-c=") {
					arg = arg[len("-c="):]
					if len(arg) > 0 {
						file = arg
					}
				} else if strings.HasPrefix(arg, "--config=") {
					arg = arg[len("--config="):]
					if len(arg) > 0 {
						file = arg
					}
				} else if strings.HasPrefix(arg, "-config=") {
					arg = arg[len("-config="):]
					if len(arg) > 0 {
						file = arg
					}
				}
			}
		}
	}
	if len(file) <= 0 {
		return cfg, nil
	}
	err := config.LoadToMap(file, cfg)
	if err != nil {
		if os.IsNotExist(err) {
			if len(cfg) <= 0 {
				h.logger.Warnf("config file %s not exist", file)
			} else {
				h.logger.Debugf("config file %s not exist", file)
			}
			return cfg, nil
		}
		h.logger.Errorf("fail to load config: %s", err)
		return nil, err
	}
	h.logger.Debugf("using config file: %s", file)
	return cfg, nil
}

func (h *Hub) loadProviders(config map[string]interface{}) error {
	h.providersMap = map[string][]*providerContext{}
	err := h.doLoadProviders(config, "providers")
	if err != nil {
		return err
	}
	list := config["providers"]
	if list != nil {
		switch providers := list.(type) {
		case []interface{}:
			for _, item := range providers {
				if cfg, ok := item.(map[string]interface{}); ok {
					err = h.addProvider("", cfg)
					if err != nil {
						return nil
					}
				} else {
					return fmt.Errorf("invalid provider config type: %v", reflect.TypeOf(cfg))
				}
			}
		case map[string]interface{}:
			err = h.doLoadProviders(providers, "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Hub) doLoadProviders(config map[string]interface{}, filter string) error {
	for key, cfg := range config {
		key = strings.ReplaceAll(key, "_", "-")
		if key == filter {
			continue
		}
		err := h.addProvider(key, cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Hub) addProvider(key string, cfg interface{}) error {
	name, label := key, ""
	idx := strings.Index(key, "@")
	if idx > 0 {
		name = key[0:idx]
		label = key[idx+1:]
	}
	if cfg != nil {
		if v, ok := cfg.(map[string]interface{}); ok {
			if val, ok := v["_name"]; ok {
				if n, ok := val.(string); ok {
					name = n
				}
			}
			if val, ok := v["_enable"]; ok {
				if enable, ok := val.(bool); ok && !enable {
					return nil
				}
			}
		}
	}
	if len(name) <= 0 {
		return fmt.Errorf("provider name must not be empty")
	}
	define, ok := serviceProviders[name]
	if !ok {
		return fmt.Errorf("provider %s not exist", name)
	}
	provider := define.Creator()()
	h.providersMap[name] = append(h.providersMap[name], &providerContext{h, key, label, name, cfg, provider, define})
	return nil
}
