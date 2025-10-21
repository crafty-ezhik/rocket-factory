package order

import (
	"context"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (r *repository) Update(_ context.Context, data serviceModel.UpdateOrderInfo, kind serviceModel.OrderUpdateKind) error {
	//r.mu.Lock()
	//order := r.data[data.UUID]
	//defer r.mu.Unlock()
	//
	//switch kind {
	//case serviceModel.OrderUpdateUPDATEINFO:
	//	order.TransactionUUID = data.UUID
	//	order.PaymentMethod = data.PaymentMethod
	//	order.Status = serviceModel.OrderStatusPAID
	//
	//case serviceModel.OrderUpdateCANCEL:
	//	order.Status = serviceModel.OrderStatusCANCELLED
	//}

	return nil
}
