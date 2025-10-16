package order

import (
	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"

	"github.com/google/uuid"
)

func (r *repository) Get(orderID uuid.UUID) (serviceModel.Order, error) {
	return serviceModel.Order{}, nil
}
