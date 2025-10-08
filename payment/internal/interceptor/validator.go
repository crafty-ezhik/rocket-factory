package interceptor

import (
	"context"
	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidatorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if err := req.(*paymentV1.PayOrderRequest).Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
		}
		return handler(ctx, req)
	}
}
