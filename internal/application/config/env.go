package config

import (
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	App struct {
		Environment string `env:"APP_ENV"     envDefault:"development"`
	}
	Debug struct {
		User string `env:"DEBUG_USER"     envDefault:"hoge"`
	}
	HTTP struct {
		Port               int      `env:"PORT"                 envDefault:"8080"`
		CORSAllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envSeparator:"," envDefault:"[]"`
	}
	Redis struct {
		Addr     string `env:"REDIS_URL"     envDefault:"0.0.0.0:6379"`
		Password string `env:"REDIS_PASS"     envDefault:"chat"`
	}
	Postgres struct {
		Host     string `env:"POSTGRES_HOST"     envDefault:"0.0.0.0"`
		Port     int    `env:"POSTGRES_PORT"     envDefault:"5432"`
		User     string `env:"POSTGRES_USER"     envDefault:"chat"`
		Pass     string `env:"POSTGRES_PASS"     envDefault:"chat"`
		Database string `env:"POSTGRES_DATABASE" envDefault:"chat"`

		ConnMaxLifetime time.Duration `env:"POSTGRES_CONN_MAX_LIFETIME" envDefault:"5m"` //  sets the maximum amount of time a connection may be reused
		MaxIdleConns    int           `env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"0"`     // sets the maximum number of connections in the idle
		MaxOpenConns    int           `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"5"`     // sets the maximum number of connections in the idle
	}
	Firebase struct {
		Credentials string `env:"FIREBASE_CREDENTIALS"     envDefault:"hoge"`
	}
}

func NewConfigFromEnv() *Config {
	var c Config

	if err := env.Parse(&c.App); err != nil {
		panic(err)
	}

	if err := env.Parse(&c.Debug); err != nil {
		panic(err)
	}

	if err := env.Parse(&c.HTTP); err != nil {
		panic(err)
	}

	if err := env.Parse(&c.Redis); err != nil {
		panic(err)
	}

	if err := env.Parse(&c.Postgres); err != nil {
		panic(err)
	}

	if err := env.Parse(&c.Firebase); err != nil {
		panic(err)
	}

	return &c
}
