package interceptors

import (
	"context"
	"errors"
	businessErrs "github.com/crafty-ezhik/rocket-factory/shared/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, convertError(err)
		}
		return resp, nil
	}
}

func convertError(err error) error {
	if businessErr := businessErrs.GetBusinessError(err); businessErr != nil {
		return businessErrs.BusinessErrorToGRPCStatus(businessErr).Err()
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return status.Error(codes.DeadlineExceeded, "request timeout exceeded")
	}
	if errors.Is(err, context.Canceled) {
		return status.Error(codes.Canceled, "request canceled by client")
	}

	// Проверка, что ошибка уже является gRPC статусом
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	// Неизвестная ошибка → Internal
	return status.Error(codes.Internal, "internal server error")
}
