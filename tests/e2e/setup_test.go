package e2e

import (
	"testing"

	"dryka.pl/trader/internal/application/config"
	"dryka.pl/trader/internal/application/server"
	"dryka.pl/trader/internal/domain/trade/model"
	"dryka.pl/trader/internal/domain/trade/provider"
	"dryka.pl/trader/internal/domain/trade/service"
	"dryka.pl/trader/internal/infrastructure/persistence/sqlite"
	"dryka.pl/trader/internal/infrastructure/persistence/sqlite/repository"
	"dryka.pl/trader/tests/database"
	"dryka.pl/trader/tests/mock"
	"github.com/caarlos0/env/v6"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	AppDependencies server.Dependencies
}

func TestE2eInfrastructureSuite(t *testing.T) {
	s := new(Suite)
	opts := env.Options{Environment: map[string]string{
		"DATABASE_FILE": "database/sqlite/database_test.sqlite",
	}}
	c, err := config.NewConfig(opts)
	if err != nil {
		t.Fatal("invalid config")
	}

	connection, err := sqlite.NewConnection("../../" + c.GetDatabaseFile())
	if err != nil {
		t.Fatal("cannot connect")
	}

	order := c.GetOrder()
	price, err := decimal.NewFromString(order.Price)
	if err != nil {
		t.Fatal("invalid price")
	}

	orderService := service.NewOrderService(model.Order{
		Quantity:       decimal.NewFromInt(order.Quantity),
		Price:          price,
		SourceCurrency: order.SourceCurrency,
		TargetCurrency: order.TargetCurrency,
	}, repository.NewAuditRepository(connection), provider.NewTimeProvider(), nil)
	s.AppDependencies = server.Dependencies{
		Logger:  mock.NewNullLogger(),
		Config:  c,
		Service: orderService,
		DB:      connection,
	}

	err = database.CreateFromTemplate(
		"../../database/sqlite/database.sqlite.template",
		"../../database/sqlite/database_test.sqlite",
	)
	if err != nil {
		t.Fatal(err)
	}

	suite.Run(t, s)
}
