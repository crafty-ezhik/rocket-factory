package grpc

import (
	"context"
	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/google/uuid"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter serviceModel.PartsFilter) ([]serviceModel.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod string) (string, error)
}
