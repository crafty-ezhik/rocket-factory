package order

import (
	"context"
	"database/sql"
	"errors"
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, orderID uuid.UUID) (serviceModel.Order, error) {
	builderSelect := sq.Select("*").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_uuid": orderID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return serviceModel.Order{}, err
	}

	row, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return serviceModel.Order{}, serviceModel.ErrOrderNotFound
		}
		return serviceModel.Order{}, err
	}
	var order repoModel.Order

	for row.Next() {
		err = row.Scan(
			&order.UUID,
			&order.PartUUIDs,
			&order.TotalPrice,
			&order.TransactionUUID,
			&order.PaymentMethod,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return serviceModel.Order{}, err
		}
	}

	return converter.OrderToServiceModel(order), nil
}
