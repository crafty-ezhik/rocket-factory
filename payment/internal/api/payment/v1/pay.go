package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	orderUUID, err := uuid.Parse(req.OrderUuid)
	if err != nil {
		return &paymentV1.PayOrderResponse{}, status.Error(codes.InvalidArgument, "order uuid is not valid")
	}

	userUUID, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return &paymentV1.PayOrderResponse{}, status.Error(codes.InvalidArgument, "user uuid is not valid")
	}

	transactionUUID, err := a.paymentService.PayOrder(ctx, orderUUID, userUUID, req.PaymentMethod.String())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, "request timeout exceeded")
		}
		if errors.Is(err, context.Canceled) {
			return nil, status.Error(codes.Canceled, "request canceled by client")
		}

		return &paymentV1.PayOrderResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
