package order

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
	"github.com/crafty-ezhik/rocket-factory/platform/pkg/logger"
)

func (r *repository) Update(ctx context.Context, order serviceModel.Order) error {
	builderUpdate := sq.Update(ordersTable).
		PlaceholderFormat(sq.Dollar).
		Set(orderFieldTotalPrice, order.TotalPrice).
		Set(orderFieldTransactionUUID, order.TransactionUUID).
		Set(orderFieldStatus, order.Status).
		Set(orderFieldPaymentMethod, order.PaymentMethod).
		Set(orderFieldUpdatedAt, time.Now()).
		Where(sq.Eq{orderFieldOrderUUID: order.UUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		logger.Error(ctx, "Ошибка при преобразовании запроса к SQL", zap.Error(err))
		return fmt.Errorf("build update query: %w", err)
	}

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "Ошибка при обновлении заказа", zap.Error(err))
		return fmt.Errorf("execute update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return serviceModel.ErrOrderNotFound
	}

	return nil
}
