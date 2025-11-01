package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/crafty-ezhik/rocket-factory/payment/internal/app"
	"github.com/crafty-ezhik/rocket-factory/payment/internal/config"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

const configPath = "../deploy/compose/payment/.env"

func main() {
	// Загружаем конфиг
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("❌ Ошибка загрузки конфига: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
