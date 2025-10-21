package v1

import (
	"context"
	"errors"
	"net/http"

	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) OrderCreate(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.OrderCreateRes, error) {
	if err := req.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "Invalid Create Request",
		}, nil
	}

	orderUUID, totalPrice, err := a.orderService.Create(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
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

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
