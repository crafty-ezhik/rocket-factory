package app

import (
	"context"
	"fmt"
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
