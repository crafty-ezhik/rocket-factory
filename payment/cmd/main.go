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

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

const (
	HOST          = "localhost"
	PathToSwagger = "./shared/pkg/swagger/payment/v1"
	grpcPort      = 50051
	httpPort      = 8081
)

// paymentService - реализует gRPC сервис для работы с оплатами
type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

// PayOrder - обрабатывает команду на оплату и возвращает transaction_uuid
func (ps *paymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUUID := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s\n", transactionUUID.String())

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}

func main() {
	// Открываем для прослушивания tcp соединение на порту grpcPort
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
		return
	}

	// По окончанию работы сервера, закрываем tcp соединение
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Fatalf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем наш сервис paymentService
	service := &paymentService{}
	paymentV1.RegisterPaymentServiceServer(grpcServer, service)

	// Включаем рефлексию для откладки, чтобы клиент мог видеть доступные методы
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

	// Запускаем HTTP сервер с gRPC gateway и Swagger UI
	var gwServer *http.Server
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Создаем мультиплексор для запросов
		mux := runtime.NewServeMux()

		// Настраиваем опции для соединения с gRPC- сервером. Отключаем защищенное соединение
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		// Регистрируем gRPC-gateway хендлеры
		err = paymentV1.RegisterPaymentServiceHandlerFromEndpoint(
			ctx,
			mux,
			fmt.Sprintf("%s:%d", HOST, grpcPort),
			opts,
		)
		if err != nil {
			log.Printf("failed to register gateway: %v\n", err)
			return
		}

		// Создаем файловый сервер для swagger-ui
		fileServer := http.FileServer(http.Dir(PathToSwagger))

		// Создаем HTTP маршрутизатор
		httpMux := http.NewServeMux()

		// Регистрируем API ручки
		httpMux.Handle("/api/v1/payment", mux)

		// Swagger UI endpoints
		httpMux.Handle("/swagger-ui.html", fileServer)
		httpMux.Handle("/payment.swagger.json", fileServer)

		// Настройка редиректа с корня на Swagger UI
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
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to serve HTTP: %v\n", err)
			return
		}
	}()

	// Реализуем Graceful shutdown
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
