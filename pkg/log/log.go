package log

import (
	"github.com/go-kit/kit/log"
	stdlog "log"
	"os"
	"strings"
)

func NewLogger(env string) (logger log.Logger)  {
	defer func() {
		stdlog.SetOutput(log.NewStdlibAdapter(logger))
	}()

	if strings.ToUpper(env) != "LOCAL" && env != "" {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		return logger
	}
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	return logger
}
