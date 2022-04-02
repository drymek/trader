package server

import (
	"database/sql"

	"dryka.pl/trader/internal/application/config"
	"dryka.pl/trader/internal/domain/trade/service"
	"dryka.pl/trader/internal/infrastructure/logger"
)

type Dependencies struct {
	Logger  logger.TraderLogger
	Config  config.Config
	Service service.OrderService
	DB      *sql.DB
}
