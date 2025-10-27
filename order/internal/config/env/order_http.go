package env

import (
	"github.com/caarlos0/env/v11"
	"net"
	"time"
)

type orderHTTPEnvConfig struct {
	Host            string        `env:"HTTP_HOST,required"`
	Port            string        `env:"HTTP_PORT,required"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT,required"`
}

type orderHTTPConfig struct {
	raw orderHTTPEnvConfig
}

func NewOrderHTTPConfig() (*orderHTTPConfig, error) {
	var raw orderHTTPEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &orderHTTPConfig{raw: raw}, nil
}

func (cfg *orderHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}

func (cfg *orderHTTPConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}

func (cfg *orderHTTPConfig) ShutdownTimeout() time.Duration {
	return cfg.raw.ShutdownTimeout
}
