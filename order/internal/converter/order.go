package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func OrderToHTTP(order model.Order) *orderV1.OrderDto {
	return &orderV1.OrderDto{
		OrderUUID:       order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUIDToHTTP(order.TransactionUUID),
		PaymentMethod:   paymentMethodToHTTP(order.PaymentMethod),
		Status:          orderStatusToHTTP(order.Status),
		CreatedAt:       createAtToHTTP(order.CreatedAt),
		UpdatedAt:       updateAtToHTTP(order.UpdatedAt),
	}
}

func transactionUUIDToHTTP(transactionUUID uuid.UUID) orderV1.OptNilUUID {
	return orderV1.OptNilUUID{
		Value: transactionUUID,
		Set:   true,
		Null:  false,
	}
}

func paymentMethodToHTTP(paymentMethod model.PaymentMethod) orderV1.OptNilPaymentMethod {
	out := orderV1.OptNilPaymentMethod{
		Set:  true,
		Null: false,
	}
	switch paymentMethod {
	case model.PaymentMethodUNKNOWN:
		out.Value = orderV1.PaymentMethodUNKNOWN
	case model.PaymentMethodCARD:
		out.Value = orderV1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		out.Value = orderV1.PaymentMethodSBP
	case model.PaymentMethodCREDITCARD:
		out.Value = orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodINVESTORMONEY:
		out.Value = orderV1.PaymentMethodINVESTORMONEY
	default:
		out.Value = orderV1.PaymentMethodUNKNOWN
	}
	return out
}

func orderStatusToHTTP(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPAID:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return orderV1.OrderStatusCANCELLED
	case model.OrderStatusASSEMBLED:
		return orderV1.OrderStatusASSEMBLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func createAtToHTTP(date time.Time) orderV1.OptDateTime {
	return orderV1.OptDateTime{Value: date, Set: true}
}

func updateAtToHTTP(date *time.Time) orderV1.OptDateTime {
	if date == nil {
		return orderV1.OptDateTime{Value: time.Time{}, Set: false}
	}
	return orderV1.OptDateTime{Value: *date, Set: true}
}
