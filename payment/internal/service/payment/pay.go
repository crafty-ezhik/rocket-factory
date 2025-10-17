package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
)

// PayOrder - обрабатывает команду на оплату и возвращает transaction_uuid
func (s *Service) PayOrder(_ context.Context, orderID, userID uuid.UUID, paymentMethod string) (string, error) {
	log.Printf(`
	💳 [Order Paid]
	• 🆔 Order UUID: %s
	• 👤 User UUID: %s
	• 💰 Payment Method: %s
	`, orderID.String(), userID.String(), paymentMethod,
	)
	transactionUUID := uuid.NewString()

	return transactionUUID, nil
}
