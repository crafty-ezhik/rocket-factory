package model

import (
	"errors"
	sharedErr "github.com/crafty-ezhik/rocket-factory/platform/pkg/grpc/errors"
)

var (
	ErrInvalidCredentials   = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid user UUID"))
	ErrSessionUUIDIsMissing = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("session UUID is missing"))
	ErrUserInfoIsMissing    = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("user info is missing"))
	ErrInvalidSessionUUID   = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid session UUID"))
	ErrInvalidUserUUID      = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid user UUID"))
	ErrInvalidEmail         = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("invalid email"))
	ErrPasswordIsRequired   = sharedErr.NewBusinessError(sharedErr.BadRequestErrCode, errors.New("password is required"))
	ErrUserNotFound         = sharedErr.NewBusinessError(sharedErr.NotFoundErrCode, errors.New("user not found"))
	ErrWeakPassword         = sharedErr.NewBusinessError(sharedErr.NotFoundErrCode, errors.New("password must be at least 8 characters"))
)
