package order

import (
	"context"

	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, orderID uuid.UUID) (serviceModel.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[orderID]
	if !ok {
		return serviceModel.Order{}, serviceModel.ErrNotFound
	}

	return converter.OrderToServiceModel(order), nil
}
