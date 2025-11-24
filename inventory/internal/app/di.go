package app

import (
	"context"
	"fmt"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	middlewareGRPC "github.com/crafty-ezhik/rocket-factory/platform/pkg/middleware/grpc"
	auth_v1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/auth/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	inventoryV1API "github.com/crafty-ezhik/rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/config"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/repository"
	inventoryRepository "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/part"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/service"
	inventoryService "github.com/crafty-ezhik/rocket-factory/inventory/internal/service/part"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API      inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.InventoryRepository
	mongoDBClient       *mongo.Client
	mongoDBHandle       *mongo.Database
	iamClient           middlewareGRPC.IAMClient
}

// NewDIContainer - возвращает пустой diContainer
func NewDIContainer() *diContainer {
	return &diContainer{}
}

// InventoryV1API - создает экземпляр api хендлеров
func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.PartService(ctx))
	}
	return d.inventoryV1API
}

// PartService - создает экземпляр сервиса
func (d *diContainer) PartService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.PartRepository(ctx))
	}
	return d.inventoryService
}

// PartRepository - создает экземпляр репозитория
func (d *diContainer) PartRepository(ctx context.Context) repository.InventoryRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = inventoryRepository.NewRepository(ctx, d.MongoDBHandle(ctx))

		// Добавление деталей в Mongo
		// d.inventoryRepository.Init()
	}
	return d.inventoryRepository
}

// MongoDBClient - создает клиента MongoDB и добавляет функцию закрытия в closer
func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err = client.Ping(pingCtx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		// Добавляем закрытие ресурсов
		closer.AddNamed("MongoDB client", func(ctx context.Context) error { return client.Disconnect(ctx) })

		d.mongoDBClient = client
	}
	return d.mongoDBClient
}

// MongoDBHandle - возвращает базу данных для работы в MongoDB
func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}
	return d.mongoDBHandle
}

func (d *diContainer) IAMClient(ctx context.Context) middlewareGRPC.IAMClient {
	if d.iamClient == nil {
		grpcIAM := auth_v1.NewAuthServiceClient(d.IAMConn(ctx))
		d.iamClient = grpcIAM
	}
	return d.iamClient
}

func (d *diContainer) IAMConn(_ context.Context) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		config.AppConfig().IamGRPC.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		panic(fmt.Sprintf("❌ Ошибка подключения к IAM Service: %v", err))
	}

	closer.AddNamed("IAM client", func(ctx context.Context) error {
		if err := conn.Close(); err != nil {
			logger.Error(ctx, "❌ Ошибка при закрытии подключения с IAM Service", zap.Error(err))
			return err
		}
		return nil
	})

	return conn
}
