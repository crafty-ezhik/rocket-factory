package v1

import (
	"context"

	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) OrderCancel(ctx context.Context, req orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error) {
	return nil, nil
}
