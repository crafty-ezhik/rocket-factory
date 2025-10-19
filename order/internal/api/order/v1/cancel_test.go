package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (s *ApiSuite) TestOrderCancelSuccess() {
	orderUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	tests := []struct {
		name        string
		param       orderV1.OrderCancelParams
		expectedRes orderV1.OrderCancelRes
		expectedErr error
		setupMock   func()
	}{
		{
			name: "success",
			param: orderV1.OrderCancelParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.OrderCancelNoContent{},
			setupMock: func() {
				s.orderService.On("Cancel", s.ctx, orderUUID).
					Return(nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			res, err := s.api.OrderCancel(s.ctx, tt.param)
			s.Require().NoError(err)
			s.Require().Equal(tt.expectedRes, res)
		})
	}
}

func (s *ApiSuite) TestOrderCancelFailure() {
	orderUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	invalidOrderUUID := "00000000-0000-0000-0000-0000000000033333"
	dbErr := errors.New("db error")

	tests := []struct {
		name        string
		param       orderV1.OrderCancelParams
		expectedRes orderV1.OrderCancelRes
		setupMock   func()
	}{
		{
			name: "invalid uuid",
			param: orderV1.OrderCancelParams{
				OrderUUID: invalidOrderUUID,
			},
			expectedRes: &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "order uuid validation error",
			},
			setupMock: func() {},
		},
		{
			name: "order not found",
			param: orderV1.OrderCancelParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			},
			setupMock: func() {
				s.orderService.On("Cancel", s.ctx, orderUUID).
					Return(model.ErrOrderNotFound).
					Once()
			},
		},
		{
			name: "order already cancelled",
			param: orderV1.OrderCancelParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: "order has already been cancelled",
			},
			setupMock: func() {
				s.orderService.On("Cancel", s.ctx, orderUUID).
					Return(model.ErrOrderIsCancel).
					Once()
			},
		},
		{
			name: "order already paid",
			param: orderV1.OrderCancelParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: "order has already been paid for",
			},
			setupMock: func() {
				s.orderService.On("Cancel", s.ctx, orderUUID).
					Return(model.ErrOrderIsPaid).
					Once()
			},
		},
		{
			name: "service timeout",
			param: orderV1.OrderCancelParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.RequestTimeoutError{
				Code:    http.StatusRequestTimeout,
				Message: "request timeout exceeded",
			},
			setupMock: func() {
				s.orderService.On("Cancel", s.ctx, orderUUID).
					Return(context.DeadlineExceeded).
					Once()
			},
		},
		{
			name: "service canceled",
			param: orderV1.OrderCancelParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "request cancelled",
			},
			setupMock: func() {
				s.orderService.On("Cancel", s.ctx, orderUUID).
					Return(context.Canceled).
					Once()
			},
		},
		{
			name: "server internal error",
			param: orderV1.OrderCancelParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "something went wrong",
			},
			setupMock: func() {
				s.orderService.On("Cancel", s.ctx, orderUUID).
					Return(dbErr).
					Once()
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			res, err := s.api.OrderCancel(s.ctx, tt.param)

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedRes, res)
		})
	}
}
