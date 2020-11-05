package kkafka

type Message struct {
	Timestamp   string `json:"timestamp"`
	Suuid       string `json:"suuid"`
	VersionCode string `json:"app_ver"`
	Channel     string `json:"channel"`
	Event       string `json:"event"`
	UserId      string `json:"user_id"`
	PackageName string `json:"account"`
}
