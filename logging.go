package main

import (
	"bytes"
	stdlog "log"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("hopper")

var format = logging.MustStringFormatter(
	"%{color}%{module} ▶ %{color:reset} %{message}",
)

func initLogging() {
	stderr := logging.NewLogBackend(os.Stderr, "", 0)
	formatted := logging.NewBackendFormatter(stderr, format)
	leveled := logging.AddModuleLevel(formatted)
	leveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(leveled)
}

type logWriter struct {
	logger *logging.Logger
}

func (w logWriter) Write(p []byte) (n int, err error) {
	buf := bytes.NewBuffer(p)
	w.logger.Debug(buf.String())
	return len(p), nil
}

func RedirectStandardLog(logName string) {
	logger := logging.MustGetLogger(logName)

	stdlog.SetOutput(&logWriter{logger})
}

func ResetStandardLog() {
	stdlog.SetOutput(os.Stderr)
}
