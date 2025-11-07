package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/config"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/interceptor"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/closer"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/grpc/health"
	sharedIns "github.com/crafty-ezhik/rocket-factory/platform/pkg/grpc/interceptors"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

const PathToSwagger = "./shared/pkg/swagger/inventory/v1"

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	httpServer  *http.Server
	listener    net.Listener
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
	errCh := make(chan error)

	go func() {
		errCh <- a.runHTTPGateway(ctx)
	}()

	go func() {
		errCh <- a.runGRPCServer(ctx)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
		a.initHTTPGateway,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
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

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return err
		}
		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			sharedIns.UnaryErrorInterceptor(),
			interceptor.ValidatorInterceptor(),
		))
	closer.AddNamed("Inventory GRPC Server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	// Регистрируем хендлеры
	inventoryV1.RegisterInventoryServiceServer(a.grpcServer, a.diContainer.InventoryV1API(ctx))

	return nil
}

func (a *App) initHTTPGateway(ctx context.Context) error {
	// Создаем мультиплексор для запросов
	mux := runtime.NewServeMux()

	// Настраиваем опции для соединения с gRPC сервером
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Регистрируем gRPC-gateway хендлеры
	err := inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
		ctx,
		mux,
		config.AppConfig().InventoryGRPC.Address(),
		opts)
	if err != nil {
		logger.Error(ctx, "❌ failed to register HTTP-gateway", zap.Error(err))
		return err
	}

	// Создаем файловый сервер для swagger-ui
	fileServer := http.FileServer(http.Dir(PathToSwagger))

	// Создаем HTTP маршрутизатор
	httpMux := http.NewServeMux()

	// Регистрируем API ручку
	httpMux.Handle("/api/v1/inventory/", mux)

	// Swagger UI ручки
	httpMux.Handle("/swagger-ui.html", fileServer)
	httpMux.Handle("/inventory.swagger.json", fileServer)

	// Настраиваем редирект с корня на Swagger UI
	httpMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
			return
		}
		fileServer.ServeHTTP(w, r)
	}))

	// Создаем HTTP сервер
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().InventoryHTTP.Address(),
		Handler:           httpMux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	closer.AddNamed("Inventory HTTP Gateway", func(ctx context.Context) error {
		if err := a.httpServer.Shutdown(ctx); err != nil {
			logger.Error(ctx, "❌ Inventory HTTP Gateway shutdown error", zap.Error(err))
			return err
		}
		return nil
	})

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC InventoryService server listening on %s", config.AppConfig().InventoryGRPC.Address()))
	logger.Info(ctx, "=== INVENTORY SERVICE STARTED ===")
	if err := a.grpcServer.Serve(a.listener); err != nil {
		return err
	}
	return nil
}

func (a *App) runHTTPGateway(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %s\n", config.AppConfig().InventoryHTTP.Address()))
	err := a.httpServer.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(ctx, "❌ Failed to serve HTTP", zap.Error(err))
		return err
	}
	return nil
}
