package log

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/log/term"
	"os"
	"testing"
)

func TestLog(t *testing.T) {
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
	l := term.NewLogger(os.Stdout, log.NewLogfmtLogger, colorFn)
	l.Log("level", "error", "foo", "bar")
	level.Error(l).Log("msg", "tests")
}
