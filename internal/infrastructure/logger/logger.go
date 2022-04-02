package logger

import (
	"os"

	kitlogger "github.com/go-kit/log"
)

type TraderLogger interface {
	Log(...interface{}) error
}

func NewLogger() kitlogger.Logger {
	w := kitlogger.NewSyncWriter(os.Stderr)
	logger := kitlogger.NewLogfmtLogger(w)
	logger = kitlogger.With(logger, "ts", kitlogger.DefaultTimestampUTC, "caller", kitlogger.DefaultCaller)

	return logger
}
