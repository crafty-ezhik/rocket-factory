package converter

import (
	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
)

func OrderToServiceModel(order repoModel.Order) serviceModel.Order {
	return serviceModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   PaymentMethodToService(order.PaymentMethod),
		Status:          OrderStatusToService(order.Status),
	}
}

func OrderToRepoModel(order serviceModel.Order) repoModel.Order {
	return repoModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   PaymentMethodToRepo(order.PaymentMethod),
		Status:          OrderStatusToRepo(order.Status),
	}
}

func UpdateOrderInfoToRepoModel(order serviceModel.UpdateOrderInfo) repoModel.UpdateOrderInfo {
	return repoModel.UpdateOrderInfo{
		UUID:            order.UUID,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   PaymentMethodToRepo(order.PaymentMethod),
		OrderStatus:     OrderStatusToRepo(order.OrderStatus),
	}
}
