package service_test

import (
	"context"
	"testing"
	"time"

	"dryka.pl/trader/internal/domain/trade/model"
	"dryka.pl/trader/internal/domain/trade/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) TestConsiderWithNoAction() {
	order := model.Order{
		Quantity:       decimal.NewFromInt(25),
		Price:          decimal.NewFromInt(42),
		SourceCurrency: "BNB",
		TargetCurrency: "USDT",
	}
	repository := AuditRepositoryMock{}
	timeProvider := TimeProviderMock{}

	service := service.NewOrderService(order, &repository, &timeProvider, fakeCancel)
	tick := model.Tick{
		UpdateId:        400900217,
		Symbol:          "BNBUSDT",
		BestBidPrice:    decimal.NewFromFloat(40.0),
		BestBidQuantity: decimal.NewFromFloat(10.0),
		BestAskPrice:    decimal.NewFromFloat(41.0),
		BestAskQuantity: decimal.NewFromFloat(10),
	}

	err := service.Consider(tick)
	s.Nil(err)
	repository.AssertNumberOfCalls(s.T(), "Log", 0)
}

func (s *ServiceSuite) TestConsiderWithAction() {
	order := model.Order{
		Quantity:       decimal.NewFromInt(25),
		SourceCurrency: "BNB",
		TargetCurrency: "USDT",
		Price:          decimal.NewFromInt(42),
	}
	repository := AuditRepositoryMock{}
	repository.On("Log", mock.Anything).Return(nil)
	timeProvider := TimeProviderMock{}
	now := time.Now()
	timeProvider.On("Now").Return(now)

	service := service.NewOrderService(order, &repository, &timeProvider, fakeCancel)
	tick := model.Tick{
		UpdateId:        400900223,
		Symbol:          "BNBUSDT",
		BestBidPrice:    decimal.NewFromFloat(42.0),
		BestBidQuantity: decimal.NewFromFloat(5.0),
		BestAskPrice:    decimal.NewFromFloat(43),
		BestAskQuantity: decimal.NewFromFloat(10.0),
	}

	err := service.Consider(tick)

	s.Nil(err)
	repository.AssertCalled(s.T(), "Log", model.AuditLog{
		SourceQuantity: "5",
		SourceCurrency: "BNB",
		TargetCurrency: "USDT",
		TargetQuantity: "42",
		Timestamp:      now,
		UpdateId:       400900223,
	})
}

func (s *ServiceSuite) TestConsiderWithActionOverLimit() {
	order := model.Order{
		Quantity:       decimal.NewFromInt(20),
		SourceCurrency: "BNB",
		TargetCurrency: "USDT",
		Price:          decimal.NewFromInt(42),
	}
	repository := AuditRepositoryMock{}
	repository.On("Log", mock.Anything).Return(nil)
	timeProvider := TimeProviderMock{}
	now := time.Now()
	timeProvider.On("Now").Return(now)

	service := service.NewOrderService(order, &repository, &timeProvider, fakeCancel)
	tick := model.Tick{
		UpdateId:        400900235,
		Symbol:          "BNBUSDT",
		BestBidPrice:    decimal.NewFromFloat(42.5),
		BestBidQuantity: decimal.NewFromFloat(30),
		BestAskPrice:    decimal.NewFromFloat(43),
		BestAskQuantity: decimal.NewFromFloat(10.0),
	}

	err := service.Consider(tick)

	s.Nil(err)
	repository.AssertCalled(s.T(), "Log", model.AuditLog{
		SourceQuantity: "20",
		SourceCurrency: "BNB",
		TargetCurrency: "USDT",
		TargetQuantity: "42.5",
		Timestamp:      now,
		UpdateId:       400900235,
	})
}

func (s *ServiceSuite) TestCancel() {
	order := model.Order{
		Quantity:       decimal.NewFromInt(20),
		SourceCurrency: "BNB",
		TargetCurrency: "USDT",
		Price:          decimal.NewFromInt(42),
	}
	repository := AuditRepositoryMock{}
	repository.On("Log", mock.Anything).Return(nil)
	timeProvider := TimeProviderMock{}
	now := time.Now()
	timeProvider.On("Now").Return(now)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	orderService := service.NewOrderService(order, &repository, &timeProvider, cancel)
	tick := model.Tick{
		UpdateId:        400900235,
		Symbol:          "BNBUSDT",
		BestBidPrice:    decimal.NewFromFloat(42.5),
		BestBidQuantity: decimal.NewFromFloat(30),
		BestAskPrice:    decimal.NewFromFloat(43),
		BestAskQuantity: decimal.NewFromFloat(10.0),
	}

	err := orderService.Consider(tick)
	s.NoError(err)

	select {
	case _, ok := <-ctx.Done():
		if ok {
			s.Fail("context should be closed")
		}
	default:
		s.Fail("context empty")
	}
}

type AuditRepositoryMock struct {
	mock.Mock
}

func (ar *AuditRepositoryMock) Log(log model.AuditLog) error {
	args := ar.Called(log)

	toReturn := args.Get(0)
	if toReturn == nil {
		return nil
	}

	return toReturn.(error)
}

type TimeProviderMock struct {
	mock.Mock
}

func (tp *TimeProviderMock) Now() time.Time {
	args := tp.Called()
	return args.Get(0).(time.Time)
}

func fakeCancel() {}
