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
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}

func OrderToRepoModel(order serviceModel.Order) repoModel.Order {
	return repoModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
		Status:          order.Status,
	}
}
