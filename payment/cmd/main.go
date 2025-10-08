package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50051

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

	// Реализуем Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("✅ Server stopped")
}
