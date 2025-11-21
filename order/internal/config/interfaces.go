package config

import (
	"time"

	"github.com/IBM/sarama"
)

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

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
