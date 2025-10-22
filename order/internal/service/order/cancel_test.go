package order

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

/*
Success:
1. Заказ успешно отменен -> nil
2.
3.

Failure:
1. Заказ не найден -> model.Order{}, model.ErrOrderNotFound
2. Заказ уже отменен ->model.ErrOrderIsCancel
3. Заказ уже оплачен, нельзя отменить -> ErrOrderIsPaid
4. Внутренняя ошибка сервера -> model.Order{}, dbErr := errors.New("db_error")
*/

func (s *ServiceSuite) TestCancelOrder() {
	dbErr := errors.New("DB error")
	orderUUID := uuid.New()

	tests := []struct {
		name        string
		orderUUID   uuid.UUID
		order       model.Order
		setupMock   func(orderID uuid.UUID, order model.Order, err error)
		expectedErr error
	}{
		{
			name:      "success",
			orderUUID: orderUUID,
			order: model.Order{
				UUID:            orderUUID,
				UserUUID:        uuid.UUID{},
				PartUUIDs:       nil,
				TotalPrice:      0,
				TransactionUUID: uuid.UUID{},
				Status:          model.OrderStatusCANCELLED,
				CreatedAt:       time.Time{},
				UpdatedAt:       nil,
			},
			expectedErr: nil,
			setupMock: func(orderID uuid.UUID, order model.Order, err error) {
				s.repo.On("Get", s.ctx, orderID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusPENDINGPAYMENT}, nil).Once()

				s.repo.On("Update", s.ctx, order).
					Return(nil).Once()
			},
		},
		{
			name:        "order not found",
			orderUUID:   orderUUID,
			expectedErr: model.ErrOrderNotFound,
			setupMock: func(orderID uuid.UUID, order model.Order, err error) {
				s.repo.On("Get", s.ctx, orderID).
					Return(model.Order{}, model.ErrOrderNotFound).Once()
			},
		},
		{
			name:        "order already cancelled",
			orderUUID:   orderUUID,
			expectedErr: model.ErrOrderIsCancel,
			setupMock: func(orderID uuid.UUID, order model.Order, err error) {
				s.repo.On("Get", s.ctx, orderID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusCANCELLED}, nil).Once()
			},
		},
		{
			name:        "paid order cannot be cancelled",
			orderUUID:   orderUUID,
			expectedErr: model.ErrOrderIsPaid,
			setupMock: func(orderID uuid.UUID, order model.Order, err error) {
				s.repo.On("Get", s.ctx, orderID).
					Return(model.Order{UUID: orderUUID, Status: model.OrderStatusPAID}, nil).Once()
			},
		},
		{
			name:        "db error",
			orderUUID:   orderUUID,
			order:       model.Order{UUID: orderUUID},
			expectedErr: dbErr,
			setupMock: func(orderID uuid.UUID, order model.Order, err error) {
				s.repo.On("Get", s.ctx, orderID).
					Return(model.Order{UUID: orderUUID}, nil).Once()

				s.repo.On("Update", s.ctx, order).
					Return(dbErr).Once()
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock(tt.orderUUID, tt.order, tt.expectedErr)

			err := s.service.Cancel(s.ctx, tt.orderUUID)

			s.Require().Equal(tt.expectedErr, err)
		})
	}
}
