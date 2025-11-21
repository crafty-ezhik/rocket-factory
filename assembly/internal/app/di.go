package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/crafty-ezhik/rocket-factory/assembly/internal/config"
	kafkaConv "github.com/crafty-ezhik/rocket-factory/assembly/internal/converter/kafka"
	"github.com/crafty-ezhik/rocket-factory/assembly/internal/converter/kafka/decoder"
	"github.com/crafty-ezhik/rocket-factory/assembly/internal/service"
	"github.com/crafty-ezhik/rocket-factory/assembly/internal/service/consumer/order_consumer"
	"github.com/crafty-ezhik/rocket-factory/assembly/internal/service/producer/order_producer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	wrapperKafka "github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	wrapperKafkaConsumer "github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka/consumer"
	wrapperKafkaProducer "github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka/producer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/crafty-ezhik/rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderConsumerService service.ConsumerService
	orderProducerService service.OrderProducerService

	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer wrapperKafka.Consumer
	orderPaidDecoder  kafkaConv.OrderPaidDecoder

	orderAssembledProducer wrapperKafka.Producer
	syncProducer           sarama.SyncProducer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderConsumerService() service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = order_consumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.OrderProducerService())
	}
	return d.orderConsumerService
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = order_producer.NewService(d.OrderAssembledProducer())
	}
	return d.orderProducerService
}

// OrderPaidConsumer - Создает consumer, слушающего событие order.paid
func (d *diContainer) OrderPaidConsumer() wrapperKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrapperKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderPaidConsumer
}

// ConsumerGroup - Создается consumer group на основе данных из конфигурации
func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Ошибка создания consumer group: %s\n", err.Error()))
		}

		// Добавляем закрытие consumerGroup
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}
	return d.consumerGroup
}

// OrderPaidDecoder - Создается декодер для входящих событий
func (d *diContainer) OrderPaidDecoder() kafkaConv.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}
	return d.orderPaidDecoder
}

// OrderAssembledProducer - создает producer который отправляет в топик, заданный в конфигурации
func (d *diContainer) OrderAssembledProducer() wrapperKafka.Producer {
	if d.orderAssembledProducer == nil {
		d.orderAssembledProducer = wrapperKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.orderAssembledProducer
}

// SyncProducer - создает базового producer с указанными брокерами
func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Ошибка создания sync producer: %s\n", err.Error()))
		}

		// Добавляем закрытие producer
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error { return p.Close() })

		d.syncProducer = p
	}
	return d.syncProducer
}
