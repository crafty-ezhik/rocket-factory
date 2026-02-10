package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (s *ApiSuite) TestCreateOrderSuccess() {
	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	partUUIDs := []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		uuid.MustParse("00000000-0000-0000-0000-000000000003"),
	}

	orderUUID := uuid.MustParse("00000000-0000-0000-0000-000000000006")
	totalPrice := 100.99

	tests := []struct {
		name        string
		req         *orderV1.CreateOrderRequest
		params      orderV1.OrderCreateParams
		expectedRes orderV1.OrderCreateRes
		setupMock   func()
	}{
		{
			name: "success",
			req: &orderV1.CreateOrderRequest{
				UserUUID:  userUUID,
				PartUuids: partUUIDs,
			},
			params: orderV1.OrderCreateParams{
				XSessionUUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			expectedRes: &orderV1.CreateOrderResponse{
				OrderUUID:  orderUUID,
				TotalPrice: totalPrice,
			},
			setupMock: func() {
				s.orderService.On("Create", s.ctx, userUUID, partUUIDs).
					Return(orderUUID, totalPrice, nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			res, err := s.api.OrderCreate(s.ctx, tt.req, tt.params)

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedRes, res)
		})
	}
}

func (s *ApiSuite) TestCreateOrderFailure() {
	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	partUUIDs := []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		uuid.MustParse("00000000-0000-0000-0000-000000000003"),
	}
	dbErr := errors.New("something went wrong")

	tests := []struct {
		name        string
		req         *orderV1.CreateOrderRequest
		params      orderV1.OrderCreateParams
		expectedRes orderV1.OrderCreateRes
		setupMock   func()
	}{
		{
			name: "invalid request",
			req: &orderV1.CreateOrderRequest{
				UserUUID:  userUUID,
				PartUuids: nil,
			},
			params: orderV1.OrderCreateParams{
				XSessionUUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			expectedRes: &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "Invalid Create Request",
			},
			setupMock: func() {},
		},
		{
			name: "service timeout",
			req: &orderV1.CreateOrderRequest{
				UserUUID:  userUUID,
				PartUuids: partUUIDs,
			},
			params: orderV1.OrderCreateParams{
				XSessionUUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			expectedRes: &orderV1.RequestTimeoutError{
				Code:    http.StatusRequestTimeout,
				Message: "request timeout exceeded",
			},
			setupMock: func() {
				s.orderService.On("Create", s.ctx, userUUID, partUUIDs).
					Return(uuid.Nil, 0.00, context.DeadlineExceeded).
					Once()
			},
		},
		{
			name: "service canceled",
			req: &orderV1.CreateOrderRequest{
				UserUUID:  userUUID,
				PartUuids: partUUIDs,
			},
			params: orderV1.OrderCreateParams{
				XSessionUUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			expectedRes: &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "request cancelled",
			},
			setupMock: func() {
				s.orderService.On("Create", s.ctx, userUUID, partUUIDs).
					Return(uuid.Nil, 0.00, context.Canceled).
					Once()
			},
		},
		{
			name: "service internal error",
			req: &orderV1.CreateOrderRequest{
				UserUUID:  userUUID,
				PartUuids: partUUIDs,
			},
			params: orderV1.OrderCreateParams{
				XSessionUUID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			expectedRes: &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "something went wrong",
			},
			setupMock: func() {
				s.orderService.On("Create", s.ctx, userUUID, partUUIDs).
					Return(uuid.Nil, 0.00, dbErr).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()
			res, err := s.api.OrderCreate(s.ctx, tt.req, tt.params)

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedRes, res)
		})
	}
}
