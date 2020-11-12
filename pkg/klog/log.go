package klog

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	jwt2 "glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	"gorm.io/gorm/logger"
)

func NewLogger(env contract.Env) (logger log.Logger) {
	if !env.IsLocal() {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		return log.With(logger)
	}
	// Color by level value
	colorFn := func(keyvals ...interface{}) term.FgBgColor {
		for i := 0; i < len(keyvals)-1; i += 2 {
			if keyvals[i] != "level" {
				continue
			}
			if value, ok := keyvals[i+1].(level.Value); ok {
				switch value.String() {
				case "debug":
					return term.FgBgColor{Fg: term.DarkGray}
				case "info":
					return term.FgBgColor{Fg: term.Gray}
				case "warn":
					return term.FgBgColor{Fg: term.Yellow}
				case "error":
					return term.FgBgColor{Fg: term.Red}
				case "crit":
					return term.FgBgColor{Fg: term.Gray, Bg: term.DarkRed}
				default:
					return term.FgBgColor{}
				}
			}
		}
		return term.FgBgColor{}
	}
	logger = term.NewLogger(os.Stdout, log.NewLogfmtLogger, colorFn)
	return log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
}

func WithContext(logger log.Logger, ctx context.Context) log.Logger {
	claim := jwt2.GetClaim(ctx)
	transport, _ := ctx.Value("transport").(string)
	requestUrl, _ := ctx.Value("request-url").(string)

	return log.With(
		logger,
		"transport", transport,
		"requestUrl", requestUrl,
		"userId", claim.UserId,
		"suuid", claim.Suuid,
	)
}

type GormLogAdapter struct {
	Logging log.Logger
}

func (g GormLogAdapter) LogMode(logLevel logger.LogLevel) logger.Interface {
	panic("Setting GORM LogMode is not allowed for kit log")
}

func (g GormLogAdapter) Info(ctx context.Context, s string, i ...interface{}) {
	level.Info(g.Logging).Log("msg", fmt.Sprintf(s, i...))
}

func (g GormLogAdapter) Warn(ctx context.Context, s string, i ...interface{}) {
	level.Warn(g.Logging).Log("msg", fmt.Sprintf(s, i...))
}

func (g GormLogAdapter) Error(ctx context.Context, s string, i ...interface{}) {
	level.Error(g.Logging).Log("msg", fmt.Sprintf(s, i...))
}

func (g GormLogAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)

	var l log.Logger
	if err == nil {
		l = level.Debug(g.Logging)
	} else {
		l = level.Warn(g.Logging)
	}
	if rows == -1 {
		l.Log("sql", sql, "duration", elapsed, "rows", "-", "err", err)
	} else {
		l.Log("sql", sql, "duration", elapsed, "rows", rows, "err", err)
	}
}

type JaegerLogAdapter struct {
	Logging log.Logger
}

func (l JaegerLogAdapter) Infof(msg string, args ...interface{}) {
	level.Info(l.Logging).Log("msg", fmt.Sprintf(msg, args...))
}

func (l JaegerLogAdapter) Error(msg string) {
	level.Error(l.Logging).Log("msg", msg)
}

func LevelFilter(levelCfg string) level.Option {
	switch levelCfg {
	case "debug":
		return level.AllowDebug()
	case "info":
		return level.AllowInfo()
	case "warn":
		return level.AllowWarn()
	case "error":
		return level.AllowError()
	default:
		return level.AllowAll()
	}
}

type KafkaLogAdapter struct {
	Logging log.Logger
}

func (k KafkaLogAdapter) Printf(s string, i ...interface{}) {
	k.Logging.Log("msg", fmt.Sprintf(s, i...))
}
