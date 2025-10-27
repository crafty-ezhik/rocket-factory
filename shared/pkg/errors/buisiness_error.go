package errors

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type ErrorCode int64

const (
	BadRequestErrCode ErrorCode = iota
	UnauthorizedErrCode
	ForbiddenErrCode
	NotFoundErrCode
	MethodNotAllowedErrCode
	RequestTimeoutErrCode
	TooManyRequestsErrCode
	InternalServiceErrCode
	ServiceUnavailableErrCode
	CanceledErrCode
)

// businessError - структура ошибки
type businessError struct {
	code ErrorCode
	err  error
}

func (b *businessError) Error() string {
	if b.err != nil {
		return b.err.Error()
	}
	return "unknown error"
}

func (b *businessError) Code() ErrorCode { return b.code }
func (b *businessError) Unwrap() error   { return b.err }

// NewBusinessError - создает новую businessError с определенным кодом и ошибкой
func NewBusinessError(code ErrorCode, err error) *businessError {
	return &businessError{code, err}
}

// GetBusinessError - возвращает businessError, если err это businessError
func GetBusinessError(err error) *businessError {
	var businessErr *businessError
	if errors.As(err, &businessErr) {
		return businessErr
	}
	return nil
}

// BusinessErrorToGRPCStatus - конвертирует businessError в gRPC status
func BusinessErrorToGRPCStatus(err *businessError) *status.Status {
	return status.New(errCodeToGRPCCode(err.Code()), err.Error())
}

func errCodeToGRPCCode(code ErrorCode) codes.Code {
	switch code {
	case BadRequestErrCode:
		return codes.InvalidArgument
	case UnauthorizedErrCode:
		return codes.Unauthenticated
	case ForbiddenErrCode:
		return codes.PermissionDenied
	case NotFoundErrCode:
		return codes.NotFound
	case MethodNotAllowedErrCode:
		return codes.PermissionDenied
	case RequestTimeoutErrCode:
		return codes.DeadlineExceeded
	case TooManyRequestsErrCode:
		return codes.ResourceExhausted
	case InternalServiceErrCode:
		return codes.Internal
	case ServiceUnavailableErrCode:
		return codes.Unavailable
	case CanceledErrCode:
		return codes.Canceled
	default:
		return codes.Unknown
	}
}

func errCodeToHTTPCode(code ErrorCode) int {
	switch code {
	case BadRequestErrCode:
		return http.StatusBadRequest
	case UnauthorizedErrCode:
		return http.StatusUnauthorized
	case ForbiddenErrCode:
		return http.StatusForbidden
	case NotFoundErrCode:
		return http.StatusNotFound
	case MethodNotAllowedErrCode:
		return http.StatusMethodNotAllowed
	case RequestTimeoutErrCode:
		return http.StatusRequestTimeout
	case TooManyRequestsErrCode:
		return http.StatusTooManyRequests
	case InternalServiceErrCode:
		return http.StatusInternalServerError
	case ServiceUnavailableErrCode:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
