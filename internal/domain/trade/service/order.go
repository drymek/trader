package service

import (
	"context"
	"sync"

	"dryka.pl/trader/internal/domain/trade/model"
	"dryka.pl/trader/internal/domain/trade/provider"
	"dryka.pl/trader/internal/domain/trade/repository"
	"github.com/shopspring/decimal"
)

type service struct {
	order      model.Order
	repository repository.Audit
	time       provider.TimeProvider
	cancel     context.CancelFunc
	mu         *sync.Mutex
}

func (s *service) Consider(tick model.Tick) error {
	err := s.order.Validate()
	if err != nil {
		return err
	}

	if tick.Symbol != s.order.Symbol() {
		return nil
	}

	if s.order.Price.LessThanOrEqual(tick.BestBidPrice) {
		transactionQuantity := decimal.Min(s.order.Quantity, tick.BestBidQuantity)
		s.mu.Lock()
		s.order.Quantity = s.order.Quantity.Sub(transactionQuantity)
		s.mu.Unlock()

		err = s.repository.Log(model.AuditLog{
			SourceCurrency: s.order.SourceCurrency,
			SourceQuantity: transactionQuantity.String(),
			TargetCurrency: s.order.TargetCurrency,
			TargetQuantity: tick.BestBidPrice.String(),
			Timestamp:      s.time.Now(),
			UpdateId:       tick.UpdateId,
		})

		if err != nil {
			return ErrUnableToProcessOrder
		}

		if s.order.Quantity.Equal(decimal.Zero) {
			s.cancel()
		}
	}

	return nil
}

type OrderService interface {
	Consider(tick model.Tick) error
}

func NewOrderService(order model.Order, repository repository.Audit, time provider.TimeProvider, cancel context.CancelFunc) OrderService {
	return &service{
		order:      order,
		repository: repository,
		time:       time,
		cancel:     cancel,
		mu:         &sync.Mutex{},
	}
}
