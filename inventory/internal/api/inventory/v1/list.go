package v1

import (
	"context"
	"errors"
	"github.com/crafty-ezhik/rocket-factory/inventory/internal/converter"
	inventoryV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filteredParts, err := a.inventoryService.List(ctx, converter.PartsFilterToServiceModel(req.Filter))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "request timeout exceeded")
		}
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, "request canceled by client")
		}
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}
	
	return &inventoryV1.ListPartsResponse{
		Parts: converter.SlicePartToProto(filteredParts),
	}, nil
}
