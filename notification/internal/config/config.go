package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/crafty-ezhik/rocket-factory/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Kafka                  KafkaConfig
	Logger                 LoggerConfig
	OrderPaidConsumer      OrderConsumerConfig
	OrderAssembledConsumer OrderConsumerConfig
	TgBot                  TelegramBotConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	kafkaConfig, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}
	loggerConfig, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}
	orderPaidConsumerConfig, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}
	orderAssembledConsumerConfig, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}
	tgBotConfig, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}
	appConfig = &config{
		Kafka:                  kafkaConfig,
		OrderPaidConsumer:      orderPaidConsumerConfig,
		OrderAssembledConsumer: orderAssembledConsumerConfig,
		TgBot:                  tgBotConfig,
		Logger:                 loggerConfig,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
