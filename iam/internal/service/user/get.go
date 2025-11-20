package user

import (
	"context"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, userUUID uuid.UUID) (model.User, error) {
	return s.userRepo.Get(ctx, userUUID)
}
