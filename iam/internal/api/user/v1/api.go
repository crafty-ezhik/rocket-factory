package v1

import (
	"github.com/crafty-ezhik/rocket-factory/iam/internal/service"
	userV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	service service.UserService
}

func NewUserAPI(service service.UserService) *api {
	return &api{service: service}
}
