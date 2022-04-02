package model_test

import (
	"testing"

	"dryka.pl/trader/internal/domain/trade/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type ModelOrderSuite struct {
	suite.Suite
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(ModelOrderSuite))
}

func (s *ModelOrderSuite) TestInvalidQuantity() {
	order := model.Order{
		Quantity:       decimal.NewFromInt(-10),
		Price:          decimal.NewFromInt(10),
		SourceCurrency: "BNB",
		TargetCurrency: "USDT",
	}
	err := order.Validate()
	s.Error(err)
	s.Equal(model.ErrInvalidOderQuantity, err)
}

func (s *ModelOrderSuite) TestInvalidSourceCurrency() {
	order := model.Order{
		Quantity:       decimal.NewFromInt(10),
		Price:          decimal.NewFromInt(10),
		SourceCurrency: "PLN",
		TargetCurrency: "USDT",
	}
	err := order.Validate()
	s.Error(err)
	s.Equal(model.ErrInvalidOderSourceCurrency, err)
}

func (s *ModelOrderSuite) TestInvalidTargetCurrency() {
	order := model.Order{
		Quantity:       decimal.NewFromInt(10),
		Price:          decimal.NewFromInt(10),
		SourceCurrency: "BNB",
		TargetCurrency: "PLN",
	}
	err := order.Validate()
	s.Error(err)
	s.Equal(model.ErrInvalidOderTargetCurrency, err)
}
