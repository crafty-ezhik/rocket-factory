package order

import (
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
	"github.com/google/uuid"
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

	var order repoModel.Order
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&order.UUID,
		&order.UserUUID,
		&order.PartUUIDs,
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return serviceModel.Order{}, serviceModel.ErrOrderNotFound
		}
		return serviceModel.Order{}, err
	}

	return converter.OrderToServiceModel(order), nil
}
