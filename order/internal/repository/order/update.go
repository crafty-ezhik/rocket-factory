package order

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	serviceModel "github.com/crafty-ezhik/rocket-factory/order/internal/model"
)

func (r *repository) Update(ctx context.Context, order serviceModel.Order) error {
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("total_price", order.TotalPrice).
		Set("transaction_uuid", order.TransactionUUID).
		Set("status", order.Status).
		Set("payment_method", order.PaymentMethod).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"order_uuid": order.UUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute update: %w", err)
	}
	return nil
}
