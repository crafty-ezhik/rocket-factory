package integration

import (
	"context"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/testcontainers"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/testcontainers/app"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/testcontainers/mongo"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/testcontainers/network"
)

const (
	// Параметры для контейнеров
	inventoryAppName    = "inventory-app"
	inventoryImageName  = "inventory-app:latest"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"

	// Значение переменных окружения
	loggerLevelValue = "info"
	startupTimeout   = 3 * time.Minute
)

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	// Шаг 1: Создаем общую Docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Fatal(ctx, "❌ Не удалось создать общую сеть", zap.Error(err))
	}

	logger.Info(ctx, "🔥 Название сети:", zap.String("networkName", generatedNetwork.Name())) // critical: Удалить потом

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)
	mongoAuthDb := getEnvWithLogging(ctx, testcontainers.MongoAuthDBKey)

	logger.Info(ctx, "🔥 Переменные окружения MongoDB", // critical: Удалить потом
		zap.String("username", mongoUsername),
		zap.String("password", mongoPassword),
		zap.String("database", mongoDatabase),
		zap.String("image", mongoImageName),
		zap.String("authdbname", mongoAuthDb),
	)

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	logger.Info(ctx, "🔥 Параметры запуска Mongo", // critical: Удалить потом
		zap.String("network", generatedNetwork.Name()),
		zap.String("containerName", testcontainers.MongoContainerName),
		zap.String("containerImage", mongoImageName),
		zap.String("database", mongoDatabase),
		zap.String("username", mongoUsername),
		zap.String("password", mongoPassword),
	)

	// Шаг 2: Запускаем контейнер с MongoDB
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
		mongo.WithAuthDB(mongoAuthDb),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "❌ Не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	// Шаг 3: Запускаем контейнер с приложением
	// projectRoot := path.GetProjectRoot()

	logger.Info(ctx, "🔥 Имя контейнера MongoDB", zap.String("containerName", generatedMongo.Config().ContainerName)) // critical: Удалить потом

	appEnv := map[string]string{
		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	// Создаем настраиваемую стратегию ожидания с увеличенным таймаутом
	waitStrategy := wait.ForLog("=== INVENTORY SERVICE STARTED ===").
		WithStartupTimeout(90 * time.Second).
		WithPollInterval(500 * time.Millisecond)

	logger.Info(ctx, "🔥 Параметры запуска контейнера с приложением", // critical: Удалить потом
		zap.String("name", inventoryAppName),
		zap.String("port", grpcPort),
		zap.String("image", inventoryImageName),
		zap.String("networkName", generatedNetwork.Name()),
		zap.Any("env", appEnv),
	)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		// app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithImage(inventoryImageName),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "❌ Не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

// getEnvWithLogging возвращает значение переменной окружения с логированием
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "⚠️ Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
