package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	"github.com/crafty-ezhik/rocket-factory/notification/internal/client/http"
	telegramClient "github.com/crafty-ezhik/rocket-factory/notification/internal/client/http/telegram"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/config"
	kafkaConv "github.com/crafty-ezhik/rocket-factory/notification/internal/converter/kafka"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/converter/kafka/decoder"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/service"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/service/consumer/order_assembled_consumer"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/service/consumer/order_paid_consumer"
	"github.com/crafty-ezhik/rocket-factory/notification/internal/service/telegram"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	wrapperKafka "github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	wrapperKafkaConsumer "github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka/consumer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

type diContainer struct {
	telegramService service.TelegramService
	telegramClient  http.TelegramClient
	telegramBot     *bot.Bot

	orderPaidConsumerService      service.OrderPaidConsumerService
	orderAssembledConsumerService service.OrderAssembledConsumerService

	consumerGroupPaid      sarama.ConsumerGroup
	consumerGroupAssembled sarama.ConsumerGroup
	orderPaidConsumer      wrapperKafka.Consumer
	orderPaidDecoder       kafkaConv.OrderPaidDecoder
	orderAssembledConsumer wrapperKafka.Consumer
	orderAssembledDecoder  kafkaConv.OrderAssembledDecoder
}

func NewDiContainer() *diContainer { return &diContainer{} }

func (d *diContainer) TelegramService() service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegram.NewService(d.TelegramClient())
	}
	return d.telegramService
}

func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = order_paid_consumer.NewService(
			d.OrderPaidDecoder(),
			d.OrderPaidConsumer(),
			d.TelegramService(),
		)
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumerService() service.OrderAssembledConsumerService {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumerService = order_assembled_consumer.NewService(
			d.OrderAssembledConsumer(),
			d.OrderAssembledDecoder(),
			d.TelegramService(),
		)
	}
	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumer() wrapperKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrapperKafkaConsumer.NewConsumer(
			d.ConsumerGroupPaid(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			// kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderPaidConsumer
}

func (d *diContainer) OrderAssembledConsumer() wrapperKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrapperKafkaConsumer.NewConsumer(
			d.ConsumerGroupAssembled(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			// kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderAssembledConsumer
}

func (d *diContainer) ConsumerGroupPaid() sarama.ConsumerGroup {
	if d.consumerGroupPaid == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Ошибка создания consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka paid consumer group", func(ctx context.Context) error {
			return d.consumerGroupPaid.Close()
		})

		d.consumerGroupPaid = consumerGroup
	}
	return d.consumerGroupPaid
}

func (d *diContainer) ConsumerGroupAssembled() sarama.ConsumerGroup {
	if d.consumerGroupAssembled == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Ошибка создания consumer group: %s\n", err.Error()))
		}

		closer.AddNamed("Kafka paid consumer group", func(ctx context.Context) error {
			return d.consumerGroupAssembled.Close()
		})
		d.consumerGroupAssembled = consumerGroup
	}
	return d.consumerGroupAssembled
}

func (d *diContainer) OrderPaidDecoder() kafkaConv.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssembledDecoder() kafkaConv.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}
	return d.orderAssembledDecoder
}

func (d *diContainer) TelegramClient() http.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot())
	}
	return d.telegramClient
}

func (d *diContainer) TelegramBot() *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TgBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}
		d.telegramBot = b
	}
	return d.telegramBot
}
