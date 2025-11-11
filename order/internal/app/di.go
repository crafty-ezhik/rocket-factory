package app

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	googleGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/crafty-ezhik/rocket-factory/order/internal/api/order/v1"
	"github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc"
	inventoryV1GRPC "github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentV1GRPC "github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/crafty-ezhik/rocket-factory/order/internal/config"
	kafkaConv "github.com/crafty-ezhik/rocket-factory/order/internal/converter/kafka"
	"github.com/crafty-ezhik/rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository"
	orderRepo "github.com/crafty-ezhik/rocket-factory/order/internal/repository/order"
	"github.com/crafty-ezhik/rocket-factory/order/internal/service"
	"github.com/crafty-ezhik/rocket-factory/order/internal/service/consumer/order_consumer"
	orderService "github.com/crafty-ezhik/rocket-factory/order/internal/service/order"
	"github.com/crafty-ezhik/rocket-factory/order/internal/service/producer/order_producer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	wrapperKafka "github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka"
	wrapperKafkaConsumer "github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka/consumer"
	wrapperKafkaProducer "github.com/crafty-ezhik/rocket-factory/platform/pkg/kafka/producer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/crafty-ezhik/rocket-factory/platform/pkg/middleware/kafka"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API           orderV1.Handler
	orderService         service.OrderService
	orderRepository      repository.OrderRepository
	orderConsumerService service.ConsumerService
	orderProducerService service.OrderProducerService

	pgConnPool *pgxpool.Pool

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	consumerGroup          sarama.ConsumerGroup
	orderAssembledConsumer wrapperKafka.Consumer

	orderAssembledDecoder kafkaConv.OrderAssembledDecoder
	syncProducer          sarama.SyncProducer
	orderPaidProducer     wrapperKafka.Producer
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.PartService(ctx))
	}
	return d.orderV1API
}

func (d *diContainer) PartService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(d.PartRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx), d.OrderProducerService())
	}
	return d.orderService
}

// Kafka producer

// OrderProducerService - Создает сервис kafka producer
func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = order_producer.NewService(d.OrderPaidProducer())
	}
	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = order_consumer.NewService(d.OrderAssembledConsumer(), d.PartRepository(ctx), d.OrderAssembledDecoder())
	}
	return d.orderConsumerService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepo.NewRepository(d.PgConnPool(ctx))
	}
	return d.orderRepository
}

func (d *diContainer) PgConnPool(ctx context.Context) *pgxpool.Pool {
	if d.pgConnPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("❌ Ошибка подключения к базе данных: %v\n", err))
		}

		// Проверка соединения
		pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
		defer pingCancel()

		err = pool.Ping(pingCtx)
		if err != nil {
			panic(fmt.Sprintf("failed to ping Postgres: %v\n", err))
		}

		// Добавляем закрытие пула в closer
		closer.AddNamed("Postgres connection pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})

		d.pgConnPool = pool
	}
	return d.pgConnPool
}

func (d *diContainer) PaymentClient(ctx context.Context) grpc.PaymentClient {
	if d.paymentClient == nil {
		gRPCPayment := paymentV1.NewPaymentServiceClient(d.PaymentConn(ctx))
		d.paymentClient = paymentV1GRPC.NewPaymentClient(gRPCPayment)
	}
	return d.paymentClient
}

func (d *diContainer) InventoryClient(ctx context.Context) grpc.InventoryClient {
	if d.inventoryClient == nil {
		gRPCInventory := inventoryV1.NewInventoryServiceClient(d.InventoryConn(ctx))
		d.inventoryClient = inventoryV1GRPC.NewInventoryClient(gRPCInventory)
	}
	return d.inventoryClient
}

func (d *diContainer) PaymentConn(_ context.Context) *googleGRPC.ClientConn {
	conn, err := googleGRPC.NewClient(
		config.AppConfig().PaymentGRPC.Address(),
		googleGRPC.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("❌ Ошибка подключения к PaymentService: %v", err))
	}

	// Добавляем в closer закрытие подключения
	closer.AddNamed("PaymentService connection", func(ctx context.Context) error {
		if err := conn.Close(); err != nil {
			logger.Error(ctx, "❌ Ошибка при закрытии подключения с PaymentService", zap.Error(err))
			return err
		}
		return nil
	})

	return conn
}

func (d *diContainer) InventoryConn(_ context.Context) *googleGRPC.ClientConn {
	conn, err := googleGRPC.NewClient(
		config.AppConfig().InventoryGRPC.Address(),
		googleGRPC.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(fmt.Sprintf("❌ Ошибка подключения к InventoryService: %v", err))
	}

	// Добавляем в closer закрытие подключения
	closer.AddNamed("InventoryService connection", func(ctx context.Context) error {
		if err := conn.Close(); err != nil {
			logger.Error(ctx, "❌ Ошибка при закрытии подключения с InventoryService", zap.Error(err))
			return err
		}
		return nil
	})

	return conn
}

// ConsumerGroup - Создается consumer group на основе данных из конфигурации
func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("❌ Ошибка создания consumer group: %s\n", err.Error()))
		}

		// Добавляем закрытие ConsumerGroup
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

// OrderAssembledConsumer - Создается consumer с определенной consumer group и списком топиков для прослушивания
func (d *diContainer) OrderAssembledConsumer() wrapperKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrapperKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

// OrderAssembledDecoder - Создается декодер для входящих событий
func (d *diContainer) OrderAssembledDecoder() kafkaConv.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssemblerDecoder()
	}
	return d.orderAssembledDecoder
}

// SyncProducer - создает базового producer с указанными брокерами
func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
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

// OrderPaidProducer - создает producer который отправляет в топик, заданный в конфигурации
func (d *diContainer) OrderPaidProducer() wrapperKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrapperKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}
	return d.orderPaidProducer
}
