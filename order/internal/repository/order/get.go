package order

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderID uuid.UUID) (serviceModel.Order, error) {
	query, args, err := buildSelectOrderQuery(orderID).ToSql()
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

func buildSelectOrderQuery(orderID uuid.UUID) sq.SelectBuilder {
	builderSelect := sq.Select(
		"order_uuid",
		"user_uuid",
		"part_uuids",
		"total_price",
		"transaction_uuid",
		"payment_method",
		"status",
		"created_at",
		"updated_at",
	).
		From("orders").
		Where(sq.Eq{"order_uuid": orderID}).
		PlaceholderFormat(sq.Dollar)

	return builderSelect
}
