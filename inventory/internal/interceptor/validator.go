package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func ValidatorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		log.Printf("🔄 Process interceptor.ValidatorInterceptor.....")
		return handler(ctx, req)
	}
}
