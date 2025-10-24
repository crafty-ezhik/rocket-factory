package v1

import (
	"context"
	"errors"

	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

func (s *APISuite) TestPayOrder() {
	validOrderUUID := uuid.New()
	validUserUUID := uuid.New()
	validTransactionUUID := uuid.New().String()

	tests := []struct {
		name             string
		req              *paymentV1.PayOrderRequest
		setupMock        func()
		expectedResponse *paymentV1.PayOrderResponse
		expectedErrMsg   string // ожидаемая gRPC ошибка
		wantServiceCall  bool
	}{
		{
			name: "Success",
			req: &paymentV1.PayOrderRequest{
				OrderUuid:     validOrderUUID.String(),
				UserUuid:      validUserUUID.String(),
				PaymentMethod: paymentV1.PaymentMethod_CARD,
			},
			setupMock: func() {
				s.paymentService.On("PayOrder", s.ctx, validOrderUUID, validUserUUID, "CARD").
					Return(validTransactionUUID, nil).Once()
			},
			expectedResponse: &paymentV1.PayOrderResponse{
				TransactionUuid: validTransactionUUID,
			},
			expectedErrMsg:  "",
			wantServiceCall: true,
		},
		{
			name: "Invalid Order UUID",
			req: &paymentV1.PayOrderRequest{
				OrderUuid:     "invalid_uuid",
				UserUuid:      validUserUUID.String(),
				PaymentMethod: paymentV1.PaymentMethod_CARD,
			},
			setupMock:        func() {},
			expectedResponse: nil,
			expectedErrMsg:   "invalid order UUID",
			wantServiceCall:  false,
		},
		{
			name: "Invalid user UUID",
			req: &paymentV1.PayOrderRequest{
				OrderUuid:     validOrderUUID.String(),
				UserUuid:      "invalid_uuid",
				PaymentMethod: paymentV1.PaymentMethod_CARD,
			},
			setupMock:        func() {},
			expectedResponse: nil,
			expectedErrMsg:   "invalid user UUID",
			wantServiceCall:  false,
		},
		{
			name: "Service timeout",
			req: &paymentV1.PayOrderRequest{
				OrderUuid:     validOrderUUID.String(),
				UserUuid:      validUserUUID.String(),
				PaymentMethod: paymentV1.PaymentMethod_CARD,
			},
			setupMock: func() {
				s.paymentService.On("PayOrder", s.ctx, validOrderUUID, validUserUUID, "CARD").
					Return("", context.DeadlineExceeded).Once()
			},
			expectedResponse: nil,
			expectedErrMsg:   "context deadline exceeded",
			wantServiceCall:  true,
		},
		{
			name: "Service canceled",
			req: &paymentV1.PayOrderRequest{
				OrderUuid:     validOrderUUID.String(),
				UserUuid:      validUserUUID.String(),
				PaymentMethod: paymentV1.PaymentMethod_CARD,
			},
			setupMock: func() {
				s.paymentService.On("PayOrder", s.ctx, validOrderUUID, validUserUUID, "CARD").
					Return("", context.Canceled).Once()
			},
			expectedResponse: nil,
			expectedErrMsg:   "context canceled",
			wantServiceCall:  true,
		},
		{
			name: "internal error",
			req: &paymentV1.PayOrderRequest{
				OrderUuid:     validOrderUUID.String(),
				UserUuid:      validUserUUID.String(),
				PaymentMethod: paymentV1.PaymentMethod_CARD,
			},
			setupMock: func() {
				err := errors.New("something went wrong")
				s.paymentService.On("PayOrder", s.ctx, validOrderUUID, validUserUUID, "CARD").
					Return("", err).Once()
			},
			expectedResponse: nil,
			expectedErrMsg:   "something went wrong",
			wantServiceCall:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			res, err := s.api.PayOrder(context.Background(), tt.req)

			if tt.expectedErrMsg != "" {
				s.Contains(err.Error(), tt.expectedErrMsg)
				s.Nil(tt.expectedResponse, res)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedResponse, res)
			}

			if tt.wantServiceCall {
				s.paymentService.AssertNotCalled(s.T(), "PayOrder")
			}
		})
	}
}
