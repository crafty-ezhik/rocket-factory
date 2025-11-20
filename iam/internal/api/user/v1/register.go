package v1

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/converter"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	userV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	if req.Info == nil || req.Info.Info == nil {
		return &userV1.RegisterResponse{}, model.ErrUserInfoIsMissing
	}
	if req.Info.Password == "" {
		return &userV1.RegisterResponse{}, model.ErrPasswordIsRequired
	}

	userUUID, err := a.service.Register(ctx, converter.UserRegInfoToModel(req.Info))
	if err != nil {
		return &userV1.RegisterResponse{}, err
	}

	return &userV1.RegisterResponse{UserUuid: userUUID.String()}, nil
}
