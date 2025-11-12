package app

import (
	"context"
	"fmt"

	"github.com/crafty-ezhik/rocket-factory/notification/internal/config"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runPaidConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer error: %w", err)
		}
	}()

	go func() {
		if err := a.runAssembledConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer error: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		// Триггерим cancel, чтобы остановить второй компонент
		cancel()
		// Дождись завершения всех задач (если есть graceful shutdown внутри)
		<-ctx.Done()
		return err
	case <-ctx.Done():
		logger.Info(ctx, "🔔 Получен сигнал завершения работы")
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJSON(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) runPaidConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer запущен")

	service := a.diContainer.OrderPaidConsumerService()
	err := service.RunConsumer(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runAssembledConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer запущен")

	service := a.diContainer.OrderAssembledConsumerService()
	err := service.RunConsumer(ctx)
	if err != nil {
		return err
	}
	return nil
}
