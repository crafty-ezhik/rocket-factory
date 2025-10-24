package v1

import (
	"context"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/converter"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filteredParts, err := a.inventoryService.List(ctx, converter.PartsFilterToServiceModel(req.Filter))
	if err != nil {
		return nil, err
	}

	return &inventoryV1.ListPartsResponse{
		Parts: converter.SlicePartToProto(filteredParts),
	}, nil
}
