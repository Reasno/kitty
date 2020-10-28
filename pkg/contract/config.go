package contract

type ConfigReader interface {
	GetString(string) string
	GetInt(string) int
	GetStringSlice(string) []string
	GetBool(string) bool
	Get(string) interface{}
	GetFloat64(string) float64
}
