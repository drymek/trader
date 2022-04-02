package mock

import "dryka.pl/trader/internal/infrastructure/logger"

type nullLogger struct {
}

func (nullLogger) Log(...interface{}) error {
	return nil
}

func NewNullLogger() logger.TraderLogger {
	return nullLogger{}
}
