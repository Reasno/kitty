package handlers

import (
	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
	logging "github.com/Reasno/kitty/pkg/log"
)

func provideLogger() log.Logger {
	return log.With(logging.NewLogger(viper.GetString("app_env")), "service", "app")
}
