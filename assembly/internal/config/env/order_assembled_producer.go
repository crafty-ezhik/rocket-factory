package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type orderAssembledProducerEnvConfig struct {
	TopicName string `env:"ORDER_ASSEMBLED_TOPIC_NAME,required"`
}

type orderAssembledProducerConfig struct {
	raw orderAssembledProducerEnvConfig
}

func NewOrderAssembledProducerConfig() (*orderAssembledProducerConfig, error) {
	var raw orderAssembledProducerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &orderAssembledProducerConfig{raw: raw}, nil
}

func (o *orderAssembledProducerConfig) Topic() string {
	return o.raw.TopicName
}

func (o *orderAssembledProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Producer.Return.Successes = true

	return config
}
