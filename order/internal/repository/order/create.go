package order

import (
	"context"

	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order serviceModel.Order) (uuid.UUID, error) {
	order.UUID = uuid.New()
	repoOrder := converter.OrderToRepoModel(order)

	r.mu.Lock()
	r.data[repoOrder.UUID] = repoOrder
	r.mu.Unlock()

	return repoOrder.UUID, nil
}
