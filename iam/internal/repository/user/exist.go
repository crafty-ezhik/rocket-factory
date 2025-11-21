package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Exist(ctx context.Context, login string) (model.User, error) {
	query, args, err := buildSelectUserExistQuery(login).ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("buildSelectUserQuery: %w", err)
	}

	var user repoModel.User
	err = r.pool.QueryRow(ctx, query, args...).Scan(&user.UUID, &user.Info.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("queryRow: %w", err)
	}

	return converter.UserToServiceModel(user), nil
}

func buildSelectUserExistQuery(login string) squirrel.SelectBuilder {
	builder := squirrel.Select(userFieldUserUUID, userFieldPassword).
		From(usersTable).
		Where(squirrel.Eq{userFieldLogin: login}).
		PlaceholderFormat(squirrel.Dollar)

	return builder
}
