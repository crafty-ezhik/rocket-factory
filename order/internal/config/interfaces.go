package config

import "time"

type PaymentGRPCConfig interface {
	Address() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type PostgresConfig interface {
	URI() string
	DBName() string
	MigrationsDir() string
}

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
	ShutdownTimeout() time.Duration
}
