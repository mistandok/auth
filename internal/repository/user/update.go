package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	serviceModel "github.com/mistandok/auth/internal/model"
	"github.com/mistandok/auth/internal/repository/user/convert"
)

// Update user in db.
func (u *Repo) Update(ctx context.Context, user *serviceModel.UserForUpdate) error {
	repoUser := convert.ToRepoUserForUpdateFromServiceUserForUpdate(user)

	queryFormat := `
		UPDATE "%s"
		SET %s
		WHERE %s=@%s
    `

	setParams := make([]string, 0)
	setFormat := "%s=@%s"
	namedArgs := make(pgx.NamedArgs)

	if repoUser.Name != nil {
		setParams = append(setParams, fmt.Sprintf(setFormat, nameColumn, nameColumn))
		namedArgs[nameColumn] = repoUser.Name
	}

	if repoUser.Role != nil {
		setParams = append(setParams, fmt.Sprintf(setFormat, roleColumn, roleColumn))
		namedArgs[roleColumn] = repoUser.Role
	}

	if repoUser.Email != nil {
		setParams = append(setParams, fmt.Sprintf(setFormat, emailColumn, emailColumn))
		namedArgs[emailColumn] = repoUser.Email
	}

	isNeedUpdate := len(namedArgs) > 0

	if !isNeedUpdate {
		return nil
	}

	namedArgs[idColumn] = repoUser.ID

	setParams = append(setParams, fmt.Sprintf(setFormat, updatedAtColumn, updatedAtColumn))
	namedArgs[updatedAtColumn] = time.Now()

	setParamsStr := strings.Join(setParams, ", ")

	query := fmt.Sprintf(queryFormat, userTable, setParamsStr, idColumn, idColumn)

	_, err := u.pool.Exec(ctx, query, namedArgs)

	return err
}
