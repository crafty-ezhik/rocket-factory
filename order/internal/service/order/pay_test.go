package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	orderId := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	userId := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	transactionUUID := uuid.MustParse("00000000-0000-0000-0000-000000000004")
	paymentMethod := model.PaymentMethodCARD

	tests := []struct {
		name           string
		orderID        uuid.UUID
		paymentMethod  model.PaymentMethod
		expectedResult uuid.UUID
		expectedErr    error
		setupMock      func()
	}{
		{
			name:           "success",
			orderID:        orderId,
			paymentMethod:  paymentMethod,
			expectedResult: transactionUUID,
			expectedErr:    nil,
			setupMock: func() {
				s.repo.On("Get", s.ctx, orderId).
					Return(model.Order{UUID: orderId, UserUUID: userId, Status: model.OrderStatusPENDINGPAYMENT}, nil).
					Once()

				s.paymentClient.On("PayOrder", mock.Anything, orderId, userId, paymentMethod).
					Return(transactionUUID.String(), nil).
					Once()

				s.repo.On("Update", s.ctx, model.UpdateOrderInfo{
					UUID:            orderId,
					TransactionUUID: transactionUUID,
					PaymentMethod:   paymentMethod,
				},
					model.OrderUpdateUPDATEINFO,
				).
					Return(nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.service.Pay(s.ctx, tt.orderID, tt.paymentMethod)

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedResult, resp)
		})
	}
}

func (s *ServiceSuite) TestPayOrderFail() {
	orderId := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	userId := uuid.MustParse("00000000-0000-0000-0000-000000000003")
	transactionUUID := uuid.MustParse("00000000-0000-0000-0000-000000000004")
	paymentMethod := model.PaymentMethodCARD

	dbErr := errors.New("db error")
	clientErr := errors.New("client error")
	uuidValidErr := errors.New("invalid UUID format")

	tests := []struct {
		name           string
		orderID        uuid.UUID
		paymentMethod  model.PaymentMethod
		expectedResult uuid.UUID
		expectedErr    error
		setupMock      func()
	}{
		{
			name:           "order not found",
			orderID:        orderId,
			paymentMethod:  paymentMethod,
			expectedResult: uuid.Nil,
			expectedErr:    model.ErrOrderNotFound,
			setupMock: func() {
				s.repo.On("Get", s.ctx, orderId).
					Return(model.Order{}, model.ErrOrderNotFound).Once()
			},
		},
		{
			name:           "order is already paid",
			orderID:        orderId,
			paymentMethod:  paymentMethod,
			expectedResult: uuid.Nil,
			expectedErr:    model.ErrOrderCannotPay,
			setupMock: func() {
				s.repo.On("Get", s.ctx, orderId).
					Return(model.Order{UUID: orderId, Status: model.OrderStatusPAID}, nil).
					Once()
			},
		},
		{
			name:           "order is cancelled",
			orderID:        orderId,
			paymentMethod:  paymentMethod,
			expectedResult: uuid.Nil,
			expectedErr:    model.ErrOrderCannotPay,
			setupMock: func() {
				s.repo.On("Get", s.ctx, orderId).
					Return(model.Order{UUID: orderId, Status: model.OrderStatusCANCELLED}, nil).
					Once()
			},
		},
		{
			name:           "payment client error",
			orderID:        orderId,
			paymentMethod:  paymentMethod,
			expectedResult: uuid.Nil,
			expectedErr:    clientErr,
			setupMock: func() {
				s.repo.On("Get", s.ctx, orderId).
					Return(model.Order{UUID: orderId, UserUUID: userId, Status: model.OrderStatusPENDINGPAYMENT}, nil).
					Once()

				s.paymentClient.On("PayOrder", mock.Anything, orderId, userId, paymentMethod).
					Return("", clientErr).
					Once()
			},
		},
		{
			name:           "uuid parse error",
			orderID:        orderId,
			paymentMethod:  paymentMethod,
			expectedResult: uuid.Nil,
			expectedErr:    uuidValidErr,
			setupMock: func() {
				s.repo.On("Get", s.ctx, orderId).
					Return(model.Order{UUID: orderId, UserUUID: userId, Status: model.OrderStatusPENDINGPAYMENT}, nil).
					Once()

				s.paymentClient.On("PayOrder", mock.Anything, orderId, userId, paymentMethod).
					Return("00000000-0000-0000-0000-00000000000333", nil).
					Once()
			},
		},
		{
			name:           "db error",
			orderID:        orderId,
			paymentMethod:  paymentMethod,
			expectedResult: uuid.Nil,
			expectedErr:    dbErr,
			setupMock: func() {
				s.repo.On("Get", s.ctx, orderId).
					Return(model.Order{UUID: orderId, UserUUID: userId, Status: model.OrderStatusPENDINGPAYMENT}, nil).
					Once()

				s.paymentClient.On("PayOrder", mock.Anything, orderId, userId, paymentMethod).
					Return(transactionUUID.String(), nil).
					Once()

				s.repo.On("Update", s.ctx, model.UpdateOrderInfo{
					UUID:            orderId,
					TransactionUUID: transactionUUID,
					PaymentMethod:   paymentMethod,
				},
					model.OrderUpdateUPDATEINFO,
				).
					Return(dbErr).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			resp, err := s.service.Pay(s.ctx, tt.orderID, tt.paymentMethod)

			s.Require().Error(err)
			s.Require().Contains(tt.expectedErr.Error(), err.Error())
			s.Require().Equal(resp, tt.expectedResult)
		})
	}
}
