package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/crafty-ezhik/rocket-factory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

const (
	ordersTable = "orders"

	orderFieldOrderUUID       = "order_uuid"
	orderFieldUserUUID        = "user_uuid"
	orderFieldPartUuids       = "part_uuids"
	orderFieldTotalPrice      = "total_price"
	orderFieldTransactionUUID = "transaction_uuid"
	orderFieldPaymentMethod   = "payment_method"
	orderFieldStatus          = "status"
	orderFieldCreatedAt       = "created_at"
	orderFieldUpdatedAt       = "updated_at"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
