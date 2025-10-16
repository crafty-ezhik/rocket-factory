package v1

import (
	"context"
	"errors"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/converter"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	partID, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid UUID")
	}

	part, err := a.inventoryService.Get(ctx, partID)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Error(codes.NotFound, "part not found")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "request timeout exceeded")
		}
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, "request canceled by client")
		}
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
