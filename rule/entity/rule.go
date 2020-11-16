package entity

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/knadh/koanf"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
	"glab.tagtic.cn/ad_gains/kitty/rule/msg"
)

type Rule struct {
	If      string   `yaml:"if"`
	Then    dto.Data `yaml:"then"`
	program *vm.Program
}

type Config struct {
	Style string      `yaml:"style"`
	Rule  interface{} `yaml:"rule"`
}

type CentralRules struct {
	Style string `yaml:"style"`
	Rule  struct {
		List []struct {
			Name     string `yaml:"name"`
			Icon     string `yaml:"icon"`
			Path     string `yaml:"path"`
			ID       string `yaml:"id"`
			Children []struct {
				Name     string        `yaml:"name"`
				Icon     string        `yaml:"icon"`
				Path     string        `yaml:"path"`
				ID       string        `yaml:"id"`
				Children []interface{} `yaml:"children"`
			} `yaml:"children"`
		} `yaml:"list"`
	} `yaml:"rule"`
}

type CentralConfig struct {
	Style string       `yaml:"style"`
	Rule  CentralRules `yaml:"rule"`
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

func NewRules(reader io.Reader, logger log.Logger) []Rule {
	var (
		b     []byte
		err   error
		rules []Rule
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

	if c.String("style") != "advanced" {
		rule := Rule{
			If: "true",
		}
		err = c.Unmarshal("rule", &rule.Then)
		rules = append(rules, rule)
	} else {
		err = c.Unmarshal("rule", &rules)
	}
	if err != nil {
		level.Warn(logger).Log("err", errors.Wrap(err, "invalid rules"))
		rules = []Rule{}
	}

	for i := range rules {
		if err := rules[i].Compile(); err != nil {
			level.Warn(logger).Log("error", errors.Wrap(err, rules[i].If+" is not valid"), "rule", fmt.Sprintf("%+v", rules[i]))
			rules[i].If = "false"
			rules[i].Compile()
		}
		rules[i].Then = convert(rules[i].Then)
	}
	return rules
}

func (r *Rule) Compile() error {
	var err error
	r.program, err = expr.Compile(r.If, expr.Env(&dto.Payload{}))
	return err
}

func ValidateRules(reader io.Reader) error {
	var tmp []Rule

	value, err := ioutil.ReadAll(reader)
	if err != nil {
		return &ErrInvalidRules{err.Error()}
	}
	c := koanf.New(".")
	err = c.Load(rawbytes.Provider(value), kyaml.Parser())
	if err != nil {
		return &ErrInvalidRules{err.Error()}
	}

	if c.String("style") != "advanced" {
		rule := Rule{
			If: "true",
		}
		err = c.Unmarshal("rules", &rule.Then)
		tmp = append(tmp, rule)
	} else {
		err = c.Unmarshal("rules", &tmp)
	}
	if err != nil {
		return &ErrInvalidRules{err.Error()}
	}

	for i := range tmp {
		if err := tmp[i].Compile(); err != nil {
			return &ErrInvalidRules{err.Error()}
		}
	}
	return nil
}

func Calculate(rules []Rule, payload *dto.Payload, logger log.Logger) (dto.Data, error) {
	for _, rule := range rules {
		output, err := expr.Run(rule.program, payload)
		if err != nil {
			return nil, errors.Wrap(err, msg.ErrorRules)
		}
		if i, ok := output.(int); ok && i == 0 {
			level.Debug(logger).Log("msg", "negative: "+rule.If)
			continue
		}
		if b, ok := output.(bool); ok && !b {
			level.Debug(logger).Log("msg", "negative: "+rule.If)
			continue
		}
		level.Debug(logger).Log("msg", "positive: "+rule.If)
		return rule.Then, nil
	}
	return dto.Data{}, nil
}
