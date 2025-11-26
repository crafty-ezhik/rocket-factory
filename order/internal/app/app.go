package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/crafty-ezhik/rocket-factory/order/internal/config"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	HTTPMiddleware "github.com/crafty-ezhik/rocket-factory/platform/pkg/middleware/http"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
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

	// Запускаем консьюмер
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- fmt.Errorf("consumer error: %w", err)
		}
	}()

	// Запускаем HTTP-сервер
	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- fmt.Errorf("http server error: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		logger.Error(ctx, "❌ Компонент завершился с ошибкой, завершение работы", zap.Error(err))
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
		a.initHTTPServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDIContainer()
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

func (a *App) initHTTPServer(ctx context.Context) error {
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		logger.Error(ctx, "❌ Ошибка создания сервера OpenAPI", zap.Error(err))
		return err
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	authMiddleware := HTTPMiddleware.NewAuthMiddleware(a.diContainer.IAMClient(ctx))

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/api/v1/orders/ping"))
	r.Use(middleware.Timeout(config.AppConfig().OrderHTTP.ReadTimeout()))
	r.Use(authMiddleware.Handle)

	// Монтируем обработчик OpenAPI к нашему серверу
	r.Mount("/", orderServer)

	// Создаем HTTP-сервер
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	// Добавляем в closer закрытие http сервера
	closer.AddNamed("Order server", func(ctx context.Context) error {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			logger.Error(ctx, "❌ Ошибка при остановке Order Server", zap.Error(err))
			return err
		}
		return nil
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP-сервер запущен на порту %s\n", config.AppConfig().OrderHTTP.Address()))
	err := a.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "❌ Ошибка запуска сервера", zap.Error(err))
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderAssembled Kafka consumer запущен")

	err := a.diContainer.OrderConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
