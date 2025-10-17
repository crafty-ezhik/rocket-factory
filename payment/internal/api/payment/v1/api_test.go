package v1

import (
	"context"
	"errors"
	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const ServiceMethod = "PayOrder"

func (s *APISuite) TestPayOrder() {
	validOrderUUID := uuid.New()
	validUserUUID := uuid.New()
	validTransactionUUID := uuid.New().String()

	tests := []struct {
		name             string
		req              *paymentV1.PayOrderRequest
		setupMock        func()
		expectedResponse *paymentV1.PayOrderResponse
		expectedErr      error // ожидаемая gRPC ошибка
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
				s.paymentService.On(ServiceMethod, s.ctx, validOrderUUID, validUserUUID, "CARD").
					Return(validTransactionUUID, nil).Once()
			},
			expectedResponse: &paymentV1.PayOrderResponse{
				TransactionUuid: validTransactionUUID,
			},
			expectedErr:     nil,
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
			expectedErr:      status.Error(codes.InvalidArgument, "order uuid is not valid"),
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
			expectedErr:      status.Error(codes.InvalidArgument, "user uuid is not valid"),
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
				s.paymentService.On(ServiceMethod, s.ctx, validOrderUUID, validUserUUID, "CARD").
					Return("", context.DeadlineExceeded).Once()
			},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.DeadlineExceeded, "request timeout exceeded"),
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
				s.paymentService.On(ServiceMethod, s.ctx, validOrderUUID, validUserUUID, "CARD").
					Return("", context.Canceled).Once()
			},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.Canceled, "request canceled by client"),
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
				err := errors.New("db error")
				s.paymentService.On(ServiceMethod, s.ctx, validOrderUUID, validUserUUID, "CARD").
					Return("", err).Once()
			},
			expectedResponse: nil,
			expectedErr:      status.Error(codes.Internal, "db error"),
			wantServiceCall:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			res, err := s.api.PayOrder(context.Background(), tt.req)

			if tt.expectedErr != nil {
				s.Equal(tt.expectedErr, err)
				s.Nil(tt.expectedResponse, res)
			} else {
				s.NoError(err)
				s.Equal(tt.expectedResponse, res)
			}

			if tt.wantServiceCall {
				s.paymentService.AssertNotCalled(s.T(), ServiceMethod)
			}
		})
	}
}
