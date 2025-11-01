package model

import (
	"errors"

	sharedErr "github.com/crafty-ezhik/rocket-factory/platform/pkg/grpc/errors"
)

var (
	ErrInvalidUserUUID  = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid user UUID"))
	ErrInvalidOrderUUID = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid order UUID"))
)
