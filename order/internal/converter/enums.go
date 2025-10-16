package converter

import (
	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	genOrderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

// PaymentMethodToService - конвертирует PaymentMethod в сервисный тип PaymentMethod
func PaymentMethodToService(v genOrderV1.NilPaymentMethod) model.PaymentMethod {
	switch v.Value {
	case genOrderV1.PaymentMethodUNKNOWN:
		return model.PaymentMethodUNKNOWN
	case genOrderV1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case genOrderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case genOrderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case genOrderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN

	}
}

func OrderStatusToService(v genOrderV1.OrderStatus) model.OrderStatus {
	switch v {
	case genOrderV1.OrderStatusPENDINGPAYMENT:
		return model.OrderStatusPENDINGPAYMENT
	case genOrderV1.OrderStatusPAID:
		return model.OrderStatusPAID
	case genOrderV1.OrderStatusCANCELLED:
		return model.OrderStatusCANCELLED
	default:
		return model.OrderStatusPENDINGPAYMENT
	}
}
