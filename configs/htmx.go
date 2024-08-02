package configs

import "encoding/json"

type HtmxConfig struct {
	DisableInheritance bool `json:"disableInheritance"`
}

func (config HtmxConfig) Json() (result string, err error) {
	bytes, err := json.Marshal(config)

	result = string(bytes)

	return
}

func NewDefaultHtmxConfig() HtmxConfig {
	return HtmxConfig{
		DisableInheritance: true,
	}
}