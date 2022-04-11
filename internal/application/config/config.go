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
	GetCrtFile() string
	GetKeyFile() string
	GetMongoURI() string
}

type config struct {
	NewRelicConfigAppName string        `env:"NEWRELIC_CONFIG_APP_NAME" envDefault:"trader"`
	NewRelicConfigLicense string        `env:"NEWRELIC_CONFIG_LICENSE" envDefault:""`
	HttpAddr              string        `env:"HTTP_ADDR" envDefault:"0.0.0.0:8080"`
	Timeout               time.Duration `env:"TIMEOUT" envDefault:"3s"`
	DatabaseFile          string        `env:"DATABASE_FILE" envDefault:"database/sqlite/database.sqlite"`
	KeyFile               string        `env:"KEY_FILE" envDefault:"./configs/certificate.key"`
	CrtFile               string        `env:"CRT_FILE" envDefault:"./configs/certificate.crt"`
	MongoURI              string        `env:"MONGO_URI" envDefault:"mongodb://root:example@localhost:27017/"`
	Order                 order
}

func (c config) GetMongoURI() string {
	return c.MongoURI
}

func (c config) GetCrtFile() string {
	return c.CrtFile
}

func (c config) GetKeyFile() string {
	return c.KeyFile
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
