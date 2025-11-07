package order

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/order/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/order/internal/repository/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

func (r *repository) Get(ctx context.Context, orderID uuid.UUID) (serviceModel.Order, error) {
	query, args, err := buildSelectOrderQuery(orderID).ToSql()
	if err != nil {
		logger.Error(ctx, "Ошибка при преобразовании запроса к SQL", zap.Error(err))
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
		logger.Error(ctx, "Ошибка при получении заказа", zap.Error(err))
		return serviceModel.Order{}, err
	}

	return converter.OrderToServiceModel(order), nil
}

func buildSelectOrderQuery(orderID uuid.UUID) sq.SelectBuilder {
	builderSelect := sq.Select(
		orderFieldOrderUUID,
		orderFieldUserUUID,
		orderFieldPartUuids,
		orderFieldTotalPrice,
		orderFieldTransactionUUID,
		orderFieldPaymentMethod,
		orderFieldStatus,
		orderFieldCreatedAt,
		orderFieldUpdatedAt,
	).
		From(ordersTable).
		Where(sq.Eq{orderFieldOrderUUID: orderID}).
		PlaceholderFormat(sq.Dollar)

	return builderSelect
}
