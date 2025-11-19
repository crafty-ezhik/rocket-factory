package v1

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	authV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/auth/v1"
	"github.com/google/uuid"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	if req.SessionUuid == "" {
		return &authV1.WhoamiResponse{}, model.ErrSessionUUIDIsMissing
	}

	sessionUUID, err := uuid.Parse(req.SessionUuid)
	if err != nil {
		return &authV1.WhoamiResponse{}, model.ErrInvalidSessionUUID
	}

	response, err := a.service.Whoami(ctx, sessionUUID)
	if err != nil {
		return &authV1.WhoamiResponse{}, err
	}
	_ = response

	// TODO: Сделать конвертер и использовать в return
	return nil, nil
}
