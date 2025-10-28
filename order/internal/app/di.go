package app

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	googleGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/crafty-ezhik/rocket-factory/order/internal/api/order/v1"
	"github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc"
	inventoryV1GRPC "github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentV1GRPC "github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/crafty-ezhik/rocket-factory/order/internal/config"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository"
	orderRepo "github.com/crafty-ezhik/rocket-factory/order/internal/repository/order"
	"github.com/crafty-ezhik/rocket-factory/order/internal/service"
	orderService "github.com/crafty-ezhik/rocket-factory/order/internal/service/order"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API      orderV1.Handler
	orderService    service.OrderService
	orderRepository repository.OrderRepository

	pgConnPool *pgxpool.Pool

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
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
		d.orderService = orderService.NewService(d.PartRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx))
	}
	return d.orderService
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
