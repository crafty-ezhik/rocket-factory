package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/converter"
	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (s *ApiSuite) TestGetOrderSuccess() {
	orderUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	partUUIDs := []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		uuid.MustParse("00000000-0000-0000-0000-000000000010"),
	}

	order := model.Order{
		UUID:            orderUUID,
		UserUUID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		PartUUIDs:       partUUIDs,
		TotalPrice:      100,
		TransactionUUID: uuid.MustParse("00000000-0000-0000-0000-000000000004"),
		PaymentMethod:   model.PaymentMethodCARD,
		Status:          model.OrderStatusPAID,
	}

	tests := []struct {
		name        string
		param       orderV1.OrderGetParams
		expectedRes orderV1.OrderGetRes
		setupMock   func()
	}{
		{
			name: "success",
			param: orderV1.OrderGetParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: converter.OrderToHTTP(order),
			setupMock: func() {
				s.orderService.On("Get", s.ctx, orderUUID).
					Return(order, nil).
					Once()
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()

			res, err := s.api.OrderGet(s.ctx, tt.param)

			s.Require().NoError(err)
			s.Require().Equal(tt.expectedRes, res)
		})
	}
}

func (s *ApiSuite) TestGetOrderFailure() {
	orderUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	invalidOrderUUID := "00000000-0000-0000-0000-000000000003444444"

	dbErr := errors.New("db error")

	tests := []struct {
		name        string
		param       orderV1.OrderGetParams
		expectedRes orderV1.OrderGetRes
		setupMock   func()
	}{
		{
			name: "invalid uuid",
			param: orderV1.OrderGetParams{
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
			param: orderV1.OrderGetParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			},
			setupMock: func() {
				s.orderService.On("Get", s.ctx, orderUUID).
					Return(model.Order{}, model.ErrOrderNotFound).
					Once()
			},
		},
		{
			name: "service timeout",
			param: orderV1.OrderGetParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.RequestTimeoutError{
				Code:    http.StatusRequestTimeout,
				Message: "request timeout exceeded",
			},
			setupMock: func() {
				s.orderService.On("Get", s.ctx, orderUUID).
					Return(model.Order{}, context.DeadlineExceeded).
					Once()
			},
		},
		{
			name: "service canceled",
			param: orderV1.OrderGetParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "request cancelled",
			},
			setupMock: func() {
				s.orderService.On("Get", s.ctx, orderUUID).
					Return(model.Order{}, context.Canceled).
					Once()
			},
		},
		{
			name: "internal server error",
			param: orderV1.OrderGetParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "something went wrong",
			},
			setupMock: func() {
				s.orderService.On("Get", s.ctx, orderUUID).
					Return(model.Order{}, dbErr).
					Once()
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock()
			res, err := s.api.OrderGet(s.ctx, tt.param)
			s.Require().NoError(err)
			s.Require().Equal(tt.expectedRes, res)
		})
	}
}
