package entity

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/knadh/koanf"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
)

type Ruler interface {
	Unmarshal(reader *koanf.Koanf) error
	Calculate(payload *dto.Payload) (dto.Data, error)
	Compile() error
	ShouldEnrich() bool
}

func NewRuler(style string) (Ruler, error) {
	switch style {
	case "advanced":
		return NewAdvancedRule(), nil
	case "basic":
		return NewBasicRule(), nil
	case "switch":
		return NewSwitchRule(), nil
	case "":
		return NewBasicRule(), nil
	default:
		return nil, fmt.Errorf("unsupported style %s", style)
	}
}

type Config struct {
	Style string  `yaml:"style"`
	Rules []Ruler `yaml:"rule"`
}

type CentralRules struct {
	Style string `yaml:"style"`
	Rule  struct {
		List []struct {
			Name     string   `yaml:"name"`
			Icon     string   `yaml:"icon"`
			Path     string   `yaml:"path"`
			Tabs     []string `yaml:"tabs"`
			ID       string   `yaml:"id"`
			Children []struct {
				Name     string        `yaml:"name"`
				Icon     string        `yaml:"icon"`
				Path     string        `yaml:"path"`
				ID       string        `yaml:"id"`
				Tabs     []string      `yaml:"tabs"`
				Children []interface{} `yaml:"child"`
			} `yaml:"child"`
		} `yaml:"list"`
	} `yaml:"rule"`
}

type ErrInvalidRules struct {
	detail string
}

func (e *ErrInvalidRules) Error() string {
	return e.detail
}

// convert Yaml在反序列化时，会把字段反序列化成map[interface{}]interface{}
// 而这个结构在序列化json时会出错。
// 通过这个函数，把map[interface{}]interface{}用递归转为
// map[string]interface{}
func convert(i interface{}) dto.Data {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i.(dto.Data)
}

func NewRules(reader io.Reader, logger log.Logger) Ruler {
	var (
		b   []byte
		err error
	)
	c := koanf.New(".")
	b, err = ioutil.ReadAll(reader)
	if err != nil {
		level.Warn(logger).Log("error", errors.Wrap(err, "reader is not valid"))
		b = []byte("{}")
	}

	err = c.Load(rawbytes.Provider(b), kyaml.Parser())
	if err != nil {
		level.Warn(logger).Log("err", errors.Wrap(err, "cannot load yaml"))
	}

	ruler, err := NewRuler(c.String("style"))
	if err != nil {
		level.Warn(logger).Log("err", errors.Wrap(err, "invalid rules"))
		ruler = NewBasicRule()
	}
	err = ruler.Unmarshal(c)
	if err != nil {
		level.Warn(logger).Log("err", errors.Wrap(err, "invalid rules"))
	}

	err = ruler.Compile()
	if err != nil {
		level.Warn(logger).Log("err", errors.Wrap(err, "invalid rules"))
	}
	return ruler
}

func ValidateRules(reader io.Reader) error {
	var tmp Ruler

	value, err := ioutil.ReadAll(reader)
	if err != nil {
		return &ErrInvalidRules{err.Error()}
	}
	c := koanf.New(".")
	err = c.Load(rawbytes.Provider(value), kyaml.Parser())
	if err != nil {
		return &ErrInvalidRules{err.Error()}
	}
	tmp, err = NewRuler(c.String("style"))
	if err != nil {
		return &ErrInvalidRules{err.Error()}
	}
	if err = tmp.Unmarshal(c); err != nil {
		return &ErrInvalidRules{err.Error()}
	}
	if err := tmp.Compile(); err != nil {
		return &ErrInvalidRules{err.Error()}
	}
	return nil
}

func Calculate(rules Ruler, payload *dto.Payload) (dto.Data, error) {
	return rules.Calculate(payload)
}
