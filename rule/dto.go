package rule

import (
	"encoding/json"
)

type Payload struct {
	Channel     string `json:"channel"`
	VersionCode string `json:"version_code"`
	Os          uint8  `json:"os"`
	Imei        string `json:"imei"`
	Idfa        string `json:"idfa"`
	Oaid        string `json:"oaid"`
	Suuid       string `json:"suuid"`
	Mac         string `json:"mac"`
	AndroidId   string `json:"android_id"`
	Ip          string `json:"ip"`
}

func (p Payload) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

type Data map[string]interface{}

type Response struct {
	Code    uint `json:"code"`
	Message uint `json:"message"`
	Data    Data `json:"data"`
}

func (p Response) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}
