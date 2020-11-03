package rule

import (
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type Rule struct {
	If      string `yaml:"if"`
	Then    Data   `yaml:"then"`
	program *vm.Program
}

type ErrInvalidRules struct {
	detail string
}

func (e ErrInvalidRules) Error() string {
	return e.detail
}

// convert Yaml在反序列化时，会把字段反序列化成map[interface{}]interface{}
// 而这个结构在序列化json时会出错。
// 通过这个函数，把map[interface{}]interface{}用递归转为
// map[string]interface{}
func convert(i interface{}) Data {
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
	return i.(Data)
}

func NewRules(reader io.Reader, logger log.Logger) []Rule {
	var (
		b     []byte
		err   error
		rules []Rule
	)
	b, err = ioutil.ReadAll(reader)
	if err != nil {
		level.Warn(logger).Log("error", errors.Wrap(err, "reader is not valid"))
		b = []byte("{}")
	}
	err = yaml.Unmarshal(b, &rules)
	if err != nil {
		level.Warn(logger).Log("error", errors.Wrap(err, "invalid rules"))
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
	r.program, err = expr.Compile(r.If, expr.Env(&Payload{}))
	return err
}

func validateRules(reader io.Reader) error {
	var (
		tmp   []Rule
		value []byte
	)
	value, err := ioutil.ReadAll(reader)
	if err != nil {
		return ErrInvalidRules{err.Error()}
	}
	err = yaml.Unmarshal(value, &tmp)
	if err != nil {
		return ErrInvalidRules{err.Error()}
	}
	for i := range tmp {
		if err := tmp[i].Compile(); err != nil {
			return ErrInvalidRules{err.Error()}
		}
	}
	return nil
}
