package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/crafty-ezhik/rocket-factory/order/internal/api/order/v1"
	inventoryV1GRPC "github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentV1GRPC "github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/crafty-ezhik/rocket-factory/order/internal/repository/order"
	orderService "github.com/crafty-ezhik/rocket-factory/order/internal/service/order"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	grpcInventoryAddr = "localhost:50052"
	grpcPaymentAddr   = "localhost:50051"
)

func main() {
	// Создаем gRPC клиента для InventoryService
	inventoryConn, err := grpc.NewClient(
		grpcInventoryAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}

	defer func() {
		if cerr := inventoryConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	gRPCInventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

	// Создаем gRPC клиента для PaymentService
	paymentConn, err := grpc.NewClient(
		grpcPaymentAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}

	defer func() {
		if cerr := paymentConn.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	gRPCPayment := paymentV1.NewPaymentServiceClient(paymentConn)

	// Создаем клиент-обёртку над InventoryService и PaymentService
	inventoryClient := inventoryV1GRPC.NewInventoryClient(gRPCInventoryClient)
	paymentClient := paymentV1GRPC.NewPaymentClient(gRPCPayment)

	// Создаем обработчик для API
	repo := orderRepo.NewRepository()
	service := orderService.NewService(repo, inventoryClient, paymentClient)
	api := orderAPI.NewAPI(service)

	// Создаем OpenAPI сервер
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
		return
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчик OpenAPI к нашему серверу
	r.Mount("/", orderServer)

	// Создаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запускаем сервер
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Реализуем Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
