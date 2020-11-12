//go:generate mockery --name=ConfigReader
package contract

type ConfigReader interface {
	String(string) string
	Int(string) int
	Strings(string) []string
	Bool(string) bool
	Get(string) interface{}
	Float64(string) float64
	Cut(string) ConfigReader
}

type Env interface {
	IsLocal() bool
	IsTesting() bool
	IsDev() bool
	IsProd() bool
	String() string
}

type AppName interface {
	String() string
}

type PackageName interface {
	String() string
}
