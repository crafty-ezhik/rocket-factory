package order

import (
	"github.com/crafty-ezhik/rocket-factory/order/internal/client/grpc"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository"
	def "github.com/crafty-ezhik/rocket-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderPaidProducer def.OrderProducerService
}

func NewService(
	orderRepo repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
	orderPaidProducer def.OrderProducerService,
) *service {
	return &service{
		orderRepo:         orderRepo,
		inventoryClient:   inventoryClient,
		paymentClient:     paymentClient,
		orderPaidProducer: orderPaidProducer,
	}
}
