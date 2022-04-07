package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"dryka.pl/trader/internal/application/config"
	"dryka.pl/trader/internal/application/server"
	"dryka.pl/trader/internal/domain/trade/model"
	"dryka.pl/trader/internal/domain/trade/provider"
	"dryka.pl/trader/internal/domain/trade/service"
	"dryka.pl/trader/internal/infrastructure/logger"
	"dryka.pl/trader/internal/infrastructure/persistence/sqlite"
	"dryka.pl/trader/internal/infrastructure/persistence/sqlite/repository"
	"github.com/shopspring/decimal"
)

func main() {
	l := logger.NewLogger()
	err := l.Log("msg", "Starting service")
	if err != nil {
		panic(err)
	}

	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	connection, err := sqlite.NewConnection(c.GetDatabaseFile())
	if err != nil {
		panic(err)
	}

	order := c.GetOrder()
	price, err := decimal.NewFromString(order.Price)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	service := service.NewOrderService(
		model.Order{
			Quantity:       decimal.NewFromInt(order.Quantity),
			Price:          price,
			SourceCurrency: order.SourceCurrency,
			TargetCurrency: order.TargetCurrency,
		},
		repository.NewAuditRepository(connection),
		provider.NewTimeProvider(),
		cancel,
	)

	dependencies := server.Dependencies{
		Logger:  l,
		Config:  c,
		Service: service,
		DB:      connection,
	}

	muxer := http.TimeoutHandler(server.NewServer(dependencies), c.GetTimeout(), "Timeout!")
	srv := http.Server{Addr: c.GetHttpAddr(), Handler: muxer}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		select {
		case <-sigint:
		case <-ctx.Done():
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()

		if err := srv.Shutdown(ctx); err != nil {
			err := l.Log("HTTP server Shutdown: %v", err)
			if err != nil {
				panic(err)
			}
		}
	}()

	err = l.Log("transport", "http", "address", c.GetHttpAddr(), "msg", "listening")
	if err != nil {
		panic(err)
	}

	if err := srv.ListenAndServeTLS(c.GetCrtFile(), c.GetKeyFile()); err != http.ErrServerClosed {
		err := l.Log("HTTP server ListenAndServe: %v", err)
		if err != nil {
			panic(err)
		}
	}
}
