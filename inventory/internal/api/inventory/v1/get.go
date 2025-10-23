package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/inventory/internal/converter"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	partID, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, model.ErrInvalidUUID
	}

	part, err := a.inventoryService.Get(ctx, partID)
	if err != nil {
		return nil, err
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
