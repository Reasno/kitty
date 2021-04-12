package kkafka

type Message struct {
	Timestamp   string `json:"timestamp"`
	Suuid       string `json:"suuid"`
	VersionCode string `json:"app_ver"`
	Channel     string `json:"channel"`
	Event       string `json:"event"`
	UserId      string `json:"user_id"`
	PackageName string `json:"pkg"`
	Account     string `json:"account"`
	AppKey      string `json:"appkey"`
	DeviceID    string `json:"device_id"`
	OAID        string `json:"oaid"`
	AndroidID   string `json:"android_id"`
	MAC         string `json:"mac"`
	IP          string `json:"ip"`
	Platform    string `json:"platform"`
}
