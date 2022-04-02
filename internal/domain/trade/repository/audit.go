package repository

import "dryka.pl/trader/internal/domain/trade/model"

type Audit interface {
	Log(log model.AuditLog) error
}
