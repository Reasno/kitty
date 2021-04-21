package dto

import (
	"testing"
	"time"

	"github.com/antonmedv/expr"
	"github.com/stretchr/testify/assert"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

func TestPayload_HoursAgo(t *testing.T) {
	p := &Payload{}
	assert.Equal(t, p.HoursAgo("2021-01-01 00:00:00"),
		int(time.Now().Sub(time.Date(
			2021,
			01,
			01,
			0,
			0,
			0,
			0,
			time.Local,
		)).Hours()))
}

func TestDMP(t *testing.T) {
	env := Dmp{pb.DmpResp{
		AdClick:    1,
		AdComplete: 0,
		AdDisplay:  0,
		AdCtrDev:   0,
		Register:   "",
		Score:      0,
		ScoreTotal: 0,
		BlackType:  0,
		Ext:        "",
		Skynet: &pb.SkyNet{
			Register:             pb.SkyNet_RiskLevelReject,
			Login:                0,
			Fission:              0,
			Browse:               0,
			Task:                 0,
			Withdraw:             0,
			Level:                0,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}}

	code := `RegisterRisk() == 3`

	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, true, output)
}
