package order

import "github.com/google/uuid"

func (s *service) Pay(orderID uuid.UUID, paymentMethod string) (uuid.UUID, error) {
	return uuid.UUID{}, nil
}
