package user

import (
	def "github.com/crafty-ezhik/rocket-factory/iam/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ def.UserRepository = (*repository)(nil)

const (
	usersTable         = "users"
	userFieldUserUUID  = "user_uuid"
	userFieldLogin     = "login"
	userFieldEmail     = "email"
	userFieldPassword  = "password"
	userFieldCreatedAt = "created_at"
	userFieldUpdatedAt = "updated_at"

	notificationMethodsTable             = "notification_methods"
	notificationMethodsFieldUserUUID     = "user_uuid"
	notificationMethodsFieldProviderName = "provider_name"
	notificationMethodsFieldTarget       = "target"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
