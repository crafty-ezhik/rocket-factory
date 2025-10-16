package converter

import (
	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
)

// PaymentMethodToService - конвертирует repository PaymentMethod в сервисный тип PaymentMethod
func PaymentMethodToService(v repoModel.PaymentMethod) serviceModel.PaymentMethod {
	switch v {
	case repoModel.PaymentMethodUNKNOWN:
		return serviceModel.PaymentMethodUNKNOWN
	case repoModel.PaymentMethodCARD:
		return serviceModel.PaymentMethodCARD
	case repoModel.PaymentMethodSBP:
		return serviceModel.PaymentMethodSBP
	case repoModel.PaymentMethodCREDITCARD:
		return serviceModel.PaymentMethodCREDITCARD
	case repoModel.PaymentMethodINVESTORMONEY:
		return serviceModel.PaymentMethodINVESTORMONEY
	default:
		return serviceModel.PaymentMethodUNKNOWN

	}
}

// OrderStatusToService - конвертирует repository OrderStatus в сервисный тип OrderStatus
func OrderStatusToService(v repoModel.OrderStatus) serviceModel.OrderStatus {
	switch v {
	case repoModel.OrderStatusPENDINGPAYMENT:
		return serviceModel.OrderStatusPENDINGPAYMENT
	case repoModel.OrderStatusPAID:
		return serviceModel.OrderStatusPAID
	case repoModel.OrderStatusCANCELLED:
		return serviceModel.OrderStatusCANCELLED
	default:
		return serviceModel.OrderStatusPENDINGPAYMENT
	}
}

// PaymentMethodToRepo - конвертирует repository PaymentMethod в сервисный тип PaymentMethod
func PaymentMethodToRepo(v serviceModel.PaymentMethod) repoModel.PaymentMethod {
	switch v {
	case serviceModel.PaymentMethodUNKNOWN:
		return repoModel.PaymentMethodUNKNOWN
	case serviceModel.PaymentMethodCARD:
		return repoModel.PaymentMethodCARD
	case serviceModel.PaymentMethodSBP:
		return repoModel.PaymentMethodSBP
	case serviceModel.PaymentMethodCREDITCARD:
		return repoModel.PaymentMethodCREDITCARD
	case serviceModel.PaymentMethodINVESTORMONEY:
		return repoModel.PaymentMethodINVESTORMONEY
	default:
		return repoModel.PaymentMethodUNKNOWN

	}
}

// OrderStatusToRepo - конвертирует repository OrderStatus в сервисный тип OrderStatus
func OrderStatusToRepo(v serviceModel.OrderStatus) repoModel.OrderStatus {
	switch v {
	case serviceModel.OrderStatusPENDINGPAYMENT:
		return repoModel.OrderStatusPENDINGPAYMENT
	case serviceModel.OrderStatusPAID:
		return repoModel.OrderStatusPAID
	case serviceModel.OrderStatusCANCELLED:
		return repoModel.OrderStatusCANCELLED
	default:
		return repoModel.OrderStatusPENDINGPAYMENT
	}
}
