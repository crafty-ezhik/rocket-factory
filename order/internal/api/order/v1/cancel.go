package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) OrderCancel(ctx context.Context, req orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error) {
	orderUUID, err := uuid.Parse(req.OrderUUID)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "order uuid validation error",
		}, nil
	}

	err = a.orderService.Cancel(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, model.ErrOrderIsPaid) {
			return &orderV1.ConflictError{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		}

		if errors.Is(err, model.ErrOrderIsCancel) {
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
	return &orderV1.OrderCancelNoContent{}, nil
}
