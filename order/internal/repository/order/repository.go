package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/crafty-ezhik/rocket-factory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
