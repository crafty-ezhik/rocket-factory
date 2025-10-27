package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/payment/internal/model"
	paymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *API) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	orderUUID, err := uuid.Parse(req.OrderUuid)
	if err != nil {
		return &paymentV1.PayOrderResponse{}, model.ErrInvalidOrderUUID
	}

	userUUID, err := uuid.Parse(req.UserUuid)
	if err != nil {
		return &paymentV1.PayOrderResponse{}, model.ErrInvalidUserUUID
	}

	transactionUUID, err := a.paymentService.PayOrder(ctx, orderUUID, userUUID, req.PaymentMethod.String())
	if err != nil {
		return &paymentV1.PayOrderResponse{}, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
