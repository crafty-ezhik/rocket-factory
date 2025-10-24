package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	inventoryV1API "github.com/crafty-ezhik/rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/interceptor"
	inventoryRepository "github.com/crafty-ezhik/rocket-factory/inventory/internal/repository/part"
	inventoryService "github.com/crafty-ezhik/rocket-factory/inventory/internal/service/part"
	sharedIns "github.com/crafty-ezhik/rocket-factory/shared/pkg/interceptors"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

const (
	HOST          = "localhost"
	PathToSwagger = "./shared/pkg/swagger/inventory/v1"
	grpcPort      = 50052
	httpPort      = 8082
)

func main() {
	// Загружаем переменные окружения
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("❌ Ошибка загрузки .env файла: %v\n", err)
		return
	}

	ctx := context.Background()
	mongoURI := os.Getenv("INVENTORY_MONGO_URI")

	// Создаем клиент MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к Mongo: %v\n", err)
		return
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("❌ Ошибка закрытия соединения с Mongo: %v\n", err)
			return
		}
	}()

	// Проверяем соединение
	if err = client.Ping(ctx, nil); err != nil {
		log.Printf("❌ MongoDB недоступна: %v\n", err)
	}

	// Получаем базу MongoDB
	db := client.Database(os.Getenv("INVENTORY_MONGO_DB"))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("failed to close listener: %v\n", err)
		}
	}()

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.LoggerInterceptor(),
			sharedIns.UnaryErrorInterceptor(),
			interceptor.ValidatorInterceptor(),
		),
	)

	// Регистрируем сервис inventoryService
	repo := inventoryRepository.NewRepository(db)
	service := inventoryService.NewService(repo)
	api := inventoryV1API.NewAPI(service)

	// Добавляем детали
	// repo.Init()

	inventoryV1.RegisterInventoryServiceServer(grpcServer, api)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	// Запускаем сервер
	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Поднимаем HTTP сервер для gRPC-gateway + Swagger UI
	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем мультиплексор для запросов
		mux := runtime.NewServeMux()

		// Настраиваем опции для соединения с gRPC сервером
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Регистрируем gRPC-gateway хендлеры
		err = inventoryV1.RegisterInventoryServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("%s:%d", HOST, grpcPort),
			opts)
		if err != nil {
			log.Printf("failed to register gateway: %v\n", err)
			return
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
		gwServer = &http.Server{
			Addr:              fmt.Sprintf(":%d", httpPort),
			Handler:           httpMux,
			ReadHeaderTimeout: 10 * time.Second,
		}

		// Запускаем HTTP сервер
		log.Printf("🌐 HTTP server with gRPC-Gateway and Swagger UI listening on %d\n", httpPort)
		err = gwServer.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")

	// Останавливаем HTTP сервер
	if gwServer != nil {
		shutdownctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gwServer.Shutdown(shutdownctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
		log.Println("✅ HTTP server stopped")
	}

	// Останавливаем gRPC сервер
	grpcServer.GracefulStop()
	log.Println("✅ gRPC Server stopped")
}
