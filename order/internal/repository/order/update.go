package order

import (
	"context"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
)

func (r *repository) Update(_ context.Context, data serviceModel.UpdateOrderInfo, kind serviceModel.OrderUpdateKind) error {
	repoData := converter.UpdateOrderInfoToRepoModel(data)

	r.mu.Lock()
	order := r.data[repoData.UUID]
	defer r.mu.Unlock()

	switch kind {
	case serviceModel.OrderUpdateUPDATEINFO:
		order.TransactionUUID = repoData.UUID
		order.PaymentMethod = repoData.PaymentMethod
		order.Status = repoModel.OrderStatusPAID

	case serviceModel.OrderUpdateCANCEL:
		order.Status = repoModel.OrderStatusCANCELLED
	}

	return nil
}
