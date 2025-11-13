package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type TelegramBotConfig interface {
	Token() string
	ChatID() int64
}

type OrderConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type KafkaConfig interface {
	Brokers() []string
}
