package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type sessionEnvConfig struct {
	TTL time.Duration
}

type sessionConfig struct {
	raw sessionEnvConfig
}

func NewSessionConfig() (*sessionConfig, error) {
	var raw sessionEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &sessionConfig{raw: raw}, nil
}

func (cfg *sessionConfig) TTL() time.Duration {
	return cfg.raw.TTL
}
