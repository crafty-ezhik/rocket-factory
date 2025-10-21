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

func (s *ApiSuite) TestPayOrderSuccess() {
	orderUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	transactionUUID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	tests := []struct {
		name        string
		req         *orderV1.PayOrderRequest
		params      orderV1.OrderPayParams
		expectedRes orderV1.OrderPayRes
		setupMock   func(paymentMethod orderV1.NilPaymentMethod)
	}{
		{
			name: "success",
			req: &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.NilPaymentMethod{Value: orderV1.PaymentMethodCREDITCARD},
			},
			params: orderV1.OrderPayParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.PayOrderResponse{
				TransactionUUID: transactionUUID,
			},
			setupMock: func(paymentMethod orderV1.NilPaymentMethod) {
				s.orderService.On("Pay", s.ctx, orderUUID, converter.PaymentMethodToService(paymentMethod)).
					Return(transactionUUID, nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock(tt.req.PaymentMethod)

			res, err := s.api.OrderPay(s.ctx, tt.req, tt.params)

			s.Require().Nil(err)
			s.Require().Equal(tt.expectedRes, res)
		})
	}
}

func (s *ApiSuite) TestPayOrderFailure() {
	orderUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	invalidOrderUUID := "00000000-0000-0000-0000-00000000000244444"

	dbErr := errors.New("db error")

	tests := []struct {
		name        string
		req         *orderV1.PayOrderRequest
		params      orderV1.OrderPayParams
		expectedRes orderV1.OrderPayRes
		setupMock   func(paymentMethod orderV1.NilPaymentMethod)
	}{
		{
			name: "invalid order uuid",
			req: &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.NilPaymentMethod{Value: orderV1.PaymentMethodCREDITCARD},
			},
			params: orderV1.OrderPayParams{
				OrderUUID: invalidOrderUUID,
			},
			expectedRes: &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "order uuid validation error",
			},
			setupMock: func(paymentMethod orderV1.NilPaymentMethod) {},
		},
		{
			name: "unknown payment method",
			req: &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.NilPaymentMethod{Value: orderV1.PaymentMethodUNKNOWN},
			},
			params: orderV1.OrderPayParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "unknown payment method",
			},
			setupMock: func(paymentMethod orderV1.NilPaymentMethod) {},
		},
		{
			name: "order not found",
			req: &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.NilPaymentMethod{Value: orderV1.PaymentMethodCREDITCARD},
			},
			params: orderV1.OrderPayParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "order not found",
			},
			setupMock: func(paymentMethod orderV1.NilPaymentMethod) {
				s.orderService.On("Pay", s.ctx, orderUUID, converter.PaymentMethodToService(paymentMethod)).
					Return(uuid.Nil, model.ErrOrderNotFound).
					Once()
			},
		},
		{
			name: "order already paid or cancel",
			req: &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.NilPaymentMethod{Value: orderV1.PaymentMethodCREDITCARD},
			},
			params: orderV1.OrderPayParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: "order has already been paid or cancelled",
			},
			setupMock: func(paymentMethod orderV1.NilPaymentMethod) {
				s.orderService.On("Pay", s.ctx, orderUUID, converter.PaymentMethodToService(paymentMethod)).
					Return(uuid.Nil, model.ErrOrderCannotPay).
					Once()
			},
		},
		{
			name: "service timeout",
			req: &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.NilPaymentMethod{Value: orderV1.PaymentMethodCREDITCARD},
			},
			params: orderV1.OrderPayParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.RequestTimeoutError{
				Code:    http.StatusRequestTimeout,
				Message: "request timeout exceeded",
			},
			setupMock: func(paymentMethod orderV1.NilPaymentMethod) {
				s.orderService.On("Pay", s.ctx, orderUUID, converter.PaymentMethodToService(paymentMethod)).
					Return(uuid.Nil, context.DeadlineExceeded).
					Once()
			},
		},
		{
			name: "service canceled",
			req: &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.NilPaymentMethod{Value: orderV1.PaymentMethodCREDITCARD},
			},
			params: orderV1.OrderPayParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "request cancelled",
			},
			setupMock: func(paymentMethod orderV1.NilPaymentMethod) {
				s.orderService.On("Pay", s.ctx, orderUUID, converter.PaymentMethodToService(paymentMethod)).
					Return(uuid.Nil, context.Canceled).
					Once()
			},
		},
		{
			name: "internal server error",
			req: &orderV1.PayOrderRequest{
				PaymentMethod: orderV1.NilPaymentMethod{Value: orderV1.PaymentMethodCREDITCARD},
			},
			params: orderV1.OrderPayParams{
				OrderUUID: orderUUID.String(),
			},
			expectedRes: &orderV1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "something went wrong",
			},
			setupMock: func(paymentMethod orderV1.NilPaymentMethod) {
				s.orderService.On("Pay", s.ctx, orderUUID, converter.PaymentMethodToService(paymentMethod)).
					Return(uuid.Nil, dbErr).
					Once()
			},
		},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setupMock(tt.req.PaymentMethod)
			res, err := s.api.OrderPay(s.ctx, tt.req, tt.params)
			s.Require().Nil(err)
			s.Require().Equal(tt.expectedRes, res)
		})
	}
}
