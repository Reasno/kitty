package config

// AutoGenerated by https://yaml.to-go.online/
// Used only for config validation
type AutoGenerated struct {
	Global struct {
		Version string `yaml:"version"`
		Env     string `yaml:"env"`
		HTTP    struct {
			Addr string `yaml:"addr"`
		} `yaml:"http"`
		Grpc struct {
			Addr string `yaml:"addr"`
		} `yaml:"grpc"`
		Security struct {
			Enable bool   `yaml:"enable"`
			Kid    string `yaml:"kid"`
			Key    string `yaml:"key"`
		} `yaml:"security"`
	} `yaml:"global"`
	App struct {
		Name  string `yaml:"name"`
		Redis struct {
			Addrs    []string `yaml:"addrs"`
			Database int      `yaml:"database"`
		} `yaml:"redis"`
		Gorm struct {
			Database string `yaml:"database"`
			Dsn      string `yaml:"dsn"`
		} `yaml:"gorm"`
		Jaeger struct {
			Sampler struct {
				Type  string `yaml:"type"`
				Param int    `yaml:"param"`
			} `yaml:"sampler"`
			Log struct {
				Enable bool `yaml:"enable"`
			} `yaml:"log"`
		} `yaml:"jaeger"`
		Sms struct {
			SendURL    string `yaml:"sendUrl"`
			BalanceURL string `yaml:"balanceUrl"`
			Username   string `yaml:"username"`
			Password   string `yaml:"password"`
			Tag        string `yaml:"tag"`
		} `yaml:"sms"`
	} `yaml:"app"`
}
