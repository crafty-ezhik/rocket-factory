package repository

import (
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

type OrderRepository interface {
	Create(order serviceModel.Order) (uuid.UUID, error)
	Get(orderID uuid.UUID) (serviceModel.Order, error)
	Update(data serviceModel.UpdateOrderInfo, kind serviceModel.OrderUpdateKind) error
}
