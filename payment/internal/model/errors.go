package model

import (
	"errors"

	sharedErr "github.com/crafty-ezhik/rocket-factory/shared/pkg/errors"
)

var (
	ErrInvalidUserUUID  = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid user UUID"))
	ErrInvalidOrderUUID = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid order UUID"))
)
