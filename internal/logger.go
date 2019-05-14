package internal

import (
	"fmt"
	"os"
)

type Logger interface {
	Warn(string)
	Warnp(Pos, string)
	Fatal(error)
	Fatalp(Pos, error)
}

func NewLogger(silent bool) Logger {
	if silent {
		return silentLogger{}
	} else {
		return logger{}
	}
}

type logger struct{}

func (logger) Warn(msg string)         { log("[warn] ", msg) }
func (logger) Warnp(p Pos, msg string) { log("[warn] ", p, msg) }
func (logger) Fatal(err error)         { log("[error]", err); os.Exit(1) }
func (logger) Fatalp(p Pos, err error) { log("[error]", p, err); os.Exit(1) }

type silentLogger struct{}

func (silentLogger) Warn(msg string)         { /* no-op */ }
func (silentLogger) Warnp(p Pos, msg string) { /* no-op */ }
func (silentLogger) Fatal(err error)         { log("[error]", err); os.Exit(1) }
func (silentLogger) Fatalp(p Pos, err error) { log("[error]", p, err); os.Exit(1) }

func log(args ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}
