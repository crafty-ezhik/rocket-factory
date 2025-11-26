package config

import "time"

type PostgresConfig interface {
	DBName() string
	URI() string
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type GRPCConfig interface {
	Address() string
	Host() string
	Port() string
}

type RedisConfig interface {
	Address() string
	Host() string
	Port() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}

type SessionConfig interface {
	TTL() time.Duration
}
