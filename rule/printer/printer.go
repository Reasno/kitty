package printer

import (
	"fmt"

	"glab.tagtic.cn/ad_gains/kitty/rule/client"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
)

type Printer struct {
	Payload dto.Payload
	Engine  *client.RuleEngine
}

func NewPrinter(payload dto.Payload) (*Printer, error) {
	engine, err := client.NewRuleEngine(client.Rule("printer-prod"))
	if err != nil {
		return nil, err
	}
	rtn := &Printer{
		Payload: payload,
		Engine:  engine,
	}
	return rtn, nil
}

func (p Printer) Sprintf(msg string, val ...interface{}) string {
	conf, err := p.Engine.Of("printer-prod").Payload(&p.Payload)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf(conf.String(msg), val...)
}
