package model

import (
	"errors"

	sharedErr "github.com/crafty-ezhik/rocket-factory/platform/pkg/grpc/errors"
)

var (
	ErrPartNotFound = sharedErr.NewBusinessError(sharedErr.NotFoundErrCode, errors.New("part not found"))
	ErrInvalidUUID  = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid UUID"))
)
