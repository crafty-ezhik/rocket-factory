package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUUID uuid.UUID) (model.WhoamiResponse, error) {
	sessionInfo, err := s.sessionRepo.Get(ctx, sessionUUID)
	if err != nil {
		return model.WhoamiResponse{}, err
	}

	user, err := s.userRepo.Get(ctx, sessionInfo.UserUUID)
	if err != nil {
		return model.WhoamiResponse{}, err
	}

	return model.WhoamiResponse{
		Session: sessionInfo,
		User:    user,
	}, nil
}
