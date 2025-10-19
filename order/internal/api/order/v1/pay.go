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

func (a *api) OrderPay(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.OrderPayParams) (orderV1.OrderPayRes, error) {
	orderUUID, err := uuid.Parse(params.OrderUUID)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "order uuid validation error",
		}, nil
	}

	transactionUUID, err := a.orderService.Pay(ctx, orderUUID, converter.PaymentMethodToService(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, model.ErrOrderCannotPay) {
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return &orderV1.RequestTimeoutError{
				Code:    http.StatusRequestTimeout,
				Message: "request timeout exceeded",
			}, nil
		}
		if errors.Is(err, context.Canceled) {
			return &orderV1.BadRequestError{
				Code:    http.StatusBadRequest,
				Message: "request cancelled",
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		}, nil
	}

	return &orderV1.PayOrderResponse{TransactionUUID: transactionUUID}, nil
}
