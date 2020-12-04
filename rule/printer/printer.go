package printer

import (
	"fmt"

	"glab.tagtic.cn/ad_gains/kitty/rule/client"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
)

type Printer struct {
	Rule    string
	Payload dto.Payload
	Engine  *client.RuleEngine
}

type Option func(printer *Printer)

func Rule(rule string) Option {
	return func(c *Printer) {
		c.Rule = rule
	}
}

func Engine(engine *client.RuleEngine) Option {
	return func(c *Printer) {
		c.Engine = engine
	}
}

func NewPrinter(payload dto.Payload, opt ...Option) (*Printer, error) {
	c := Printer{
		Rule:    "printer-prod",
		Engine:  nil,
		Payload: payload,
	}
	for _, f := range opt {
		f(&c)
	}
	if c.Engine == nil {
		var err error
		c.Engine, err = client.NewRuleEngine(client.Rule(c.Rule))
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

func (p Printer) Sprintf(msg string, val ...interface{}) string {
	conf, err := p.Engine.Of(p.Rule).Payload(&p.Payload)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf(conf.String(msg), val...)
}
