package v1

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/converter"
	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) OrderGet(ctx context.Context, req orderV1.OrderGetParams) (orderV1.OrderGetRes, error) {
	orderUUID, err := uuid.Parse(req.OrderUUID)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    http.StatusBadRequest,
			Message: "order uuid validation error",
		}, nil
	}

	order, err := a.orderService.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    http.StatusNotFound,
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
		log.Println(err)
		return &orderV1.InternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong",
		}, nil
	}

	return converter.OrderToHTTP(order), nil
}
