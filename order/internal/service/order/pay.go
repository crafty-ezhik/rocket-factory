package order

import (
	"context"
	"time"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, orderID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
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
	strTransactionUUID, err := s.paymentClient.PayOrder(ctxReq, order.UUID, order.UserUUID, paymentMethod)
	if err != nil {
		//logger.Error(ctx, "Превышено время запроса к InventoryService", zap.Error(err))
		return uuid.Nil, context.DeadlineExceeded
	}

	transactionUUID, err := uuid.Parse(strTransactionUUID)
	if err != nil {
		return uuid.Nil, err
	}

	// Обновляем данные по заказу
	order.Status = model.OrderStatusPAID
	order.PaymentMethod = paymentMethod
	order.TransactionUUID = transactionUUID

	err = s.orderRepo.Update(ctx, order)
	if err != nil {
		return uuid.Nil, err
	}
	return transactionUUID, nil
}
