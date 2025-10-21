package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order serviceModel.Order) (uuid.UUID, error) {
	repoOrder := converter.OrderToRepoModel(order)

	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("part_uuids", "total_price").
		Values(repoOrder.PartUUIDs, repoOrder.TotalPrice).
		Suffix("RETURNING order_uuid")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return uuid.Nil, err
	}

	var orderUUID uuid.UUID
	err = r.pool.QueryRow(ctx, query, args...).Scan(&orderUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return orderUUID, nil
}
