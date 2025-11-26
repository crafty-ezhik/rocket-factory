package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/iam/internal/converter"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	userV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/user/v1"
)

func (a *api) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	userUUID, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return nil, model.ErrInvalidUserUUID
	}

	user, err := a.service.Get(ctx, userUUID)
	if err != nil {
		return nil, err
	}

	return &userV1.GetUserResponse{
		User: converter.UserToProto(user),
	}, nil
}
