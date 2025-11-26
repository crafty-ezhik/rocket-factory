package v1

import (
	"context"

	clientConverter "github.com/crafty-ezhik/rocket-factory/order/internal/client/converter"
	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/middleware/grpc"
	generatedInventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter serviceModel.PartsFilter) ([]serviceModel.Part, error) {
	parts, err := c.generatedClient.ListParts(grpc.ForwardSessionUUIDToGRPC(ctx), &generatedInventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}
	return clientConverter.PartListToServiceModel(parts.Parts), nil
}
