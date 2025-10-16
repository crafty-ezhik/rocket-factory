package order

import (
	"github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository"
)

type service struct {
	orderRepo repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	orderRepo repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
