package user

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	"github.com/crafty-ezhik/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/crafty-ezhik/rocket-factory/iam/internal/repository/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Create(ctx context.Context, info model.UserRegistrationInfo, hashedPassword string) (uuid.UUID, error) {
	repoUser := converter.UserRegInfoToRepoModel(info)

	err := pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		usersStmt, args, err := usersInsertBuilder(repoUser, hashedPassword).ToSql()
		if err != nil {
			return fmt.Errorf("build user insert: %w", err)
		}

		var userUUID uuid.UUID
		err = tx.QueryRow(ctx, usersStmt, args...).Scan(&userUUID)
		if err != nil {
			return fmt.Errorf("insert user: %w", err)
		}

		methodsStmt, args, err := notificationInsertBuilder(userUUID, repoUser.Info.NotificationMethods).ToSql()
		_, err = tx.Exec(ctx, methodsStmt, args...)
		if err != nil {
			return fmt.Errorf("insert notification method: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Nil, nil
}

func usersInsertBuilder(info repoModel.UserRegistrationInfo, hashedPassword string) squirrel.InsertBuilder {
	return squirrel.Insert(usersTable).
		Columns(userFieldLogin, userFieldPassword, userFieldEmail).
		Values(info.Info.Login, hashedPassword, info.Info.Email).
		Suffix(fmt.Sprintf("RETURNING %s", userFieldUserUUID)).
		PlaceholderFormat(squirrel.Dollar)
}

func notificationInsertBuilder(userUUID uuid.UUID, methods []repoModel.NotificationMethod) squirrel.InsertBuilder {
	builder := squirrel.Insert(notificationMethodsTable).
		Columns(notificationMethodsFieldUserUUID, notificationMethodsFieldProviderName, notificationMethodsFieldTarget)

	for _, m := range methods {
		builder.Values(userUUID, m.ProviderName, m.Target)
	}
	return builder.PlaceholderFormat(squirrel.Dollar)
}
