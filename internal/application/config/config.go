package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config interface {
	GetNewRelicConfigAppName() string
	GetNewRelicConfigLicense() string
	GetHttpAddr() string
	GetTimeout() time.Duration
	GetDatabaseFile() string
	GetOrder() order
}

type config struct {
	NewRelicConfigAppName string        `env:"NEWRELIC_CONFIG_APP_NAME" envDefault:"trader"`
	NewRelicConfigLicense string        `env:"NEWRELIC_CONFIG_LICENSE" envDefault:""`
	HttpAddr              string        `env:"HTTP_ADDR" envDefault:"0.0.0.0:8080"`
	Timeout               time.Duration `env:"TIMEOUT" envDefault:"3s"`
	DatabaseFile          string        `env:"DATABASE_FILE" envDefault:"database/sqlite/database.sqlite"`
	Order                 order
}

func (c config) GetNewRelicConfigAppName() string {
	return c.NewRelicConfigAppName
}

func (c config) GetNewRelicConfigLicense() string {
	return c.NewRelicConfigLicense
}

func (c config) GetHttpAddr() string {
	return c.HttpAddr
}

func (c config) GetTimeout() time.Duration {
	return c.Timeout
}

func (c config) GetDatabaseFile() string {
	return c.DatabaseFile
}

func (c config) GetOrder() order {
	return c.Order
}

func NewConfig(opts ...env.Options) (Config, error) {
	c := config{}
	if err := env.Parse(&c, opts...); err != nil {
		return c, fmt.Errorf("cannot parse main config: %w", err)
	}

	o := order{}
	if err := env.Parse(&o, opts...); err != nil {
		return c, fmt.Errorf("cannot parse order config: %w", err)
	}

	c.Order = o

	return c, nil
}
