package order

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

func (s *service) Create(ctx context.Context, userID uuid.UUID, partsIDs []uuid.UUID) (uuid.UUID, float64, error) {
	partStrUUIDs := convertUUIDStoStrings(partsIDs)

	ctxReq, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	parts, err := s.inventoryClient.ListParts(ctxReq, model.PartsFilter{UUIDs: partStrUUIDs})
	if err != nil {
		logger.Error(ctx, "Превышено время запроса к InventoryService", zap.Error(err))
		return uuid.Nil, 0, context.DeadlineExceeded
	}

	totalPrice, err := checkPartsAndCountTotalPrice(partStrUUIDs, parts)
	if err != nil {
		return uuid.Nil, 0, err
	}

	newOrder := model.Order{
		UserUUID:        userID,
		PartUUIDs:       partsIDs,
		TotalPrice:      totalPrice,
		TransactionUUID: uuid.Nil,
		PaymentMethod:   model.PaymentMethodUNKNOWN,
		Status:          model.OrderStatusPENDINGPAYMENT,
	}

	orderUUID, err := s.orderRepo.Create(ctx, newOrder)
	if err != nil {
		return uuid.Nil, 0, err
	}
	return orderUUID, totalPrice, nil
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
			return 0, fmt.Errorf("%w: part with uuid %s not found", model.ErrOrderPartNotFound, UUID)
		}
	}
	return totalPrice, nil
}
