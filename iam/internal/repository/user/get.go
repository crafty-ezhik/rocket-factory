package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	serviceModel "github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, userUUID uuid.UUID) (serviceModel.User, error) {
	var user repoModel.User

	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		// Получаем пользователя
		userQuery, args, err := buildSelectUserQuery(userUUID).ToSql()
		if err != nil {
			return fmt.Errorf("build user select: %w", err)
		}

		err = tx.QueryRow(ctx, userQuery, args...).Scan(
			&user.UUID,
			&user.Info.Login,
			&user.Info.PasswordHash,
			&user.Info.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return serviceModel.ErrUserNotFound
			}
			return fmt.Errorf("query user: %w", err)
		}

		// Получаем методы уведомлений
		methodsQuery, args, err := buildSelectMethodsQuery(userUUID).ToSql()
		if err != nil {
			return fmt.Errorf("build user methods select: %w", err)
		}

		rows, err := tx.Query(ctx, methodsQuery, args...)
		if err != nil {
			return fmt.Errorf("query notification methods: %w", err)
		}

		methods, err := pgx.CollectRows[repoModel.NotificationMethod](rows, pgx.RowToStructByName[repoModel.NotificationMethod])
		if err != nil {
			return fmt.Errorf("collect notification methods: %w", err)
		}

		user.Info.NotificationMethods = methods

		return nil
	})
	if err != nil {
		return serviceModel.User{}, fmt.Errorf("get user: %w", err)
	}

	return converter.UserToServiceModel(user), nil
}

func buildSelectUserQuery(userUUID uuid.UUID) squirrel.SelectBuilder {
	builder := squirrel.Select(
		userFieldUserUUID,
		userFieldLogin,
		userFieldPassword,
		userFieldEmail,
		userFieldCreatedAt,
		userFieldUpdatedAt,
	).
		From(usersTable).
		Where(squirrel.Eq{userFieldUserUUID: userUUID}).
		PlaceholderFormat(squirrel.Dollar)

	return builder
}

func buildSelectMethodsQuery(userUUID uuid.UUID) squirrel.SelectBuilder {
	builder := squirrel.Select(
		notificationMethodsFieldProviderName,
		notificationMethodsFieldTarget,
	).From(notificationMethodsTable).
		Where(squirrel.Eq{notificationMethodsFieldUserUUID: userUUID}).
		PlaceholderFormat(squirrel.Dollar)

	return builder
}
