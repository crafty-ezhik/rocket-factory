package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (s *service) Create(ctx context.Context, userID uuid.UUID, partsIDs []uuid.UUID) (uuid.UUID, error) {
	partStrUUIDs := convertUUIDStoStrings(partsIDs)

	ctxReq, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	parts, err := s.inventoryClient.ListParts(ctxReq, model.PartsFilter{UUIDs: partStrUUIDs})
	if err != nil {
		return uuid.Nil, err
	}

	totalPrice, err := checkPartsAndCountTotalPrice(partStrUUIDs, parts)
	if err != nil {
		return uuid.Nil, err
	}

	newOrder := model.Order{
		UUID:            uuid.New(),
		UserUUID:        userID,
		PartUUIDs:       partsIDs,
		TotalPrice:      totalPrice,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   model.PaymentMethodUNKNOWN,
		Status:          model.OrderStatusPENDINGPAYMENT,
	}

	orderUUID, err := s.orderRepo.Create(ctx, newOrder)
	if err != nil {
		return uuid.Nil, err
	}
	return orderUUID, nil
}

func convertUUIDStoStrings(parts []uuid.UUID) []string {
	partsUUID := make([]string, 0, len(parts))
	for _, partUUID := range parts {
		partsUUID = append(partsUUID, partUUID.String())
	}
	return partsUUID
}

func checkPartsAndCountTotalPrice(partsUUID []string, parts []model.Part) (float64, error) {
	totalPrice := 0.0
	for _, UUID := range partsUUID {
		exist := false
		for _, part := range parts {
			if part.UUID.String() == UUID {
				exist = true
				totalPrice += part.Price
				break
			}
		}
		if !exist {
			return 0, fmt.Errorf("part with uuid %s not found", UUID)
		}
	}
	return totalPrice, nil
}
