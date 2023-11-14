package cherryProfile

import (
	"time"

	cfacade "github.com/cherry-game/cherry/facade"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
)

type (
	Config struct {
		m map[string]any
	}
)

func Wrap(val map[string]any) *Config {
	return &Config{
		m: val,
	}
}
func (p *Config) GetConfig(path ...interface{}) cfacade.ProfileCfg {
	if len(path) > 0 {
		k := path[0]
		key := cast.ToString(k)
		if vInterface, ok := p.m[key]; ok {
			tmp := make(map[string]interface{})
			switch v := vInterface.(type) {
			case []interface{}:
				for i, value := range v {
					tmp[cast.ToString(i)] = value
				}
			case map[string]interface{}:
				tmp = v
			}
			return &Config{
				m: tmp,
			}
		}
	}
	return nil

}
func (p *Config) GetString(path interface{}, defaultVal ...string) string {
	key := path.(string)
	if v, ok := p.m[key]; ok {
		return cast.ToString(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return ""

}

func (p *Config) GetBool(path interface{}, defaultVal ...bool) bool {
	key := path.(string)
	if v, ok := p.m[key]; ok {
		return cast.ToBool(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return false
}

func (p *Config) GetInt(path interface{}, defaultVal ...int) int {
	key := path.(string)
	if v, ok := p.m[key]; ok {
		return cast.ToInt(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

func (p *Config) GetInt32(path interface{}, defaultVal ...int32) int32 {
	key := path.(string)
	if v, ok := p.m[key]; ok {
		return cast.ToInt32(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

func (p *Config) GetInt64(path interface{}, defaultVal ...int64) int64 {
	key := path.(string)
	if v, ok := p.m[key]; ok {
		return cast.ToInt64(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

func (p *Config) GetDuration(path interface{}, defaultVal ...time.Duration) time.Duration {
	key := path.(string)
	if v, ok := p.m[key]; ok {
		return cast.ToDuration(v)
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return 0
}

func (p *Config) Unmarshal(value interface{}) error {
	return mapstructure.Decode(p.m, value)
}
func (p *Config) Keys() []string {
	keys := make([]string, 0, len(p.m))
	for k := range p.m {
		keys = append(keys, k)
	}
	return keys
}
func (p *Config) Size() int {
	return len(p.m)
}
