// Manages hopper logs
package main

import (
	"bytes"
	stdlog "log"
	"os"

	"github.com/op/go-logging"
)

// main hopper logger
var log = logging.MustGetLogger("hopper")

// format for hopper logs
var format = logging.MustStringFormatter(
	"%{color}%{module} ▶ %{color:reset} %{message}",
)

// Initializes logging framework for hopper
func initLogging() {
	stderr := logging.NewLogBackend(os.Stderr, "", 0)
	formatted := logging.NewBackendFormatter(stderr, format)
	leveled := logging.AddModuleLevel(formatted)
	leveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(leveled)
}

// Hopper log writer
type logWriter struct {
	logger *logging.Logger
}

// Writes data to hopper log
func (w logWriter) Write(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(p)
	w.logger.Debug(buf.String())
	return len(p), nil
}

// RedirectStandardLog switch standard library log output to hopper logs
func RedirectStandardLog(logName string) {
	logger := logging.MustGetLogger(logName)

	stdlog.SetOutput(&logWriter{logger})
}

// ResetStandardLog switch back standard library log output to stderr
func ResetStandardLog() {
	stdlog.SetOutput(os.Stderr)
}
