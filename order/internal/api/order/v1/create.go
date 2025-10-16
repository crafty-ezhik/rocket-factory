package v1

import (
	"context"
	orderV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) OrderCreate(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.OrderCreateRes, error) {
	return nil, nil
}
