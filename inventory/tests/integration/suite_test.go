package integration

import (
	"context"
	"fmt"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testsTimeout = 5 * time.Minute

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Inventory Service Integration Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("❌ Не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Загружаем .env и устанавливаем переменные в окружение
	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env"))
	if err != nil {
		logger.Fatal(suiteCtx, "❌ Не удалось загрузить .env файл", zap.Error(err))
	}

	// Устанавливаем переменные в окружение процесса
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	logger.Info(suiteCtx, "🚀 Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = ginkgo.AfterSuite(func() {
	logger.Info(context.Background(), "📉 Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
