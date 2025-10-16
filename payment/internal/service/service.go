package service

import (
	"context"

	"github.com/google/uuid"
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderID, userID uuid.UUID, paymentMethod string) (string, error)
}
