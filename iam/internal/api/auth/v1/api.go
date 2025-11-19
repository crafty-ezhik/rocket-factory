package v1

import (
	"github.com/crafty-ezhik/rocket-factory/iam/internal/service"
	authV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	service service.AuthService
}

func NewAuthAPI(service service.AuthService) *api {
	return &api{service: service}
}
