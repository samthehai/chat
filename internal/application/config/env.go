package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	App struct {
		Environment string `env:"APP_ENV"     envDefault:"development"`
	}
	HTTP struct {
		Port int `env:"HTTP_PORT" envDefault:"8080"`
	}
	Redis struct {
		Addr     string `env:"REDIS_URL"     envDefault:"0.0.0.0:6379"`
		Password string `env:"REDIS_PASS"     envDefault:"chat"`
	}
}

func NewConfigFromEnv() *Config {
	var c Config

	if err := env.Parse(&c.App); err != nil {
		panic(err)
	}

	if err := env.Parse(&c.HTTP); err != nil {
		panic(err)
	}

	if err := env.Parse(&c.Redis); err != nil {
		panic(err)
	}

	return &c
}
