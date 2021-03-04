package entity

import (
	"github.com/antonmedv/expr/vm"
	"github.com/knadh/koanf"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
)

type BasicRule struct {
	style   string
	data    dto.Data `yaml:"then"`
	program *vm.Program
}

func NewBasicRule() *BasicRule {
	return &BasicRule{style: "basic", data: dto.Data{}}
}

func (br *BasicRule) ShouldEnrich() bool {
	return false
}

func (br *BasicRule) Unmarshal(reader *koanf.Koanf) error {
	br.style = reader.String("style")
	err := reader.Unmarshal("rule", &br.data)
	if err != nil {
		return err
	}
	return nil
}

func (br *BasicRule) Compile() error {
	br.data = convert(br.data)
	return nil
}

func (br *BasicRule) Calculate(payload *dto.Payload) (dto.Data, error) {
	return br.data, nil
}
