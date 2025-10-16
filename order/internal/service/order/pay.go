package order

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderID, userID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		return uuid.Nil, err
	}

	if order.Status != model.OrderStatusPENDINGPAYMENT {
		return uuid.Nil, model.ErrOrderCannotPay
	}

	ctxReq, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	// Оплачиваем заказ
	strTransactionUUID, err := s.paymentClient.PayOrder(ctxReq, orderID, userID, paymentMethod)
	if err != nil {
		return uuid.Nil, err
	}

	transactionUUID, err := uuid.Parse(strTransactionUUID)
	if err != nil {
		return uuid.Nil, err
	}

	// Обновляем данные по заказу
	orderInfo := model.UpdateOrderInfo{
		UUID:            orderID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
	}
	err = s.orderRepo.Update(ctx, orderInfo, model.OrderUpdateUPDATEINFO)
	if err != nil {
		return uuid.Nil, err
	}
	return transactionUUID, nil
}
