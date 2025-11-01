package order

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

func (r *repository) Create(ctx context.Context, order serviceModel.Order) (uuid.UUID, error) {
	repoOrder := converter.OrderToRepoModel(order)

	builderInsert := sq.Insert(ordersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(orderFieldUserUUID, orderFieldPartUuids, orderFieldTotalPrice).
		Values(repoOrder.UserUUID, repoOrder.PartUUIDs, repoOrder.TotalPrice).
		Suffix(fmt.Sprintf("RETURNING %s", orderFieldOrderUUID))

	query, args, err := builderInsert.ToSql()
	if err != nil {
		logger.Error(ctx, "Ошибка при преобразовании запроса к SQL", zap.Error(err))
		return uuid.Nil, err
	}

	var orderUUID uuid.UUID
	err = r.pool.QueryRow(ctx, query, args...).Scan(&orderUUID)
	if err != nil {
		logger.Error(ctx, "Ошибка создании заказа", zap.Error(err))
		return uuid.Nil, err
	}

	return orderUUID, nil
}
