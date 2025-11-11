package config

import (
	"github.com/crafty-ezhik/rocket-factory/assembly/internal/config/env"
	"github.com/joho/godotenv"
	"os"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderAssembledProducer OrderAssembledConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledProducerConfig, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerConfig, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerConfig,
		Kafka:                  kafkaConfig,
		OrderAssembledProducer: orderAssembledProducerConfig,
		OrderPaidConsumer:      orderPaidConsumerConfig,
	}
	return nil
}

func AppConfig() *config { return appConfig }
