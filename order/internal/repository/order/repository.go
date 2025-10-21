package order

import (
	def "github.com/crafty-ezhik/rocket-factory/order/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
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
