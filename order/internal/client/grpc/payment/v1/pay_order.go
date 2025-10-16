package v1

import (
	"context"

	"github.com/google/uuid"

	genPaymentV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID uuid.UUID, paymentMethod string) (string, error) {
	transactionUUIDstr, err := c.generatedClient.PayOrder(ctx, &genPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: genPaymentV1.PaymentMethod(genPaymentV1.PaymentMethod_value[paymentMethod]),
	})
	if err != nil {
		return "", err
	}
	return transactionUUIDstr.TransactionUuid, nil
}
