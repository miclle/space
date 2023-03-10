package accounts

import (
	"context"

	"github.com/fox-gonic/fox/database"
	"github.com/miclle/space/accounts/params"
	"github.com/miclle/space/models"
)

// Service for account interface
type Service interface {
	CreateAccount(ctx context.Context, params *params.CreateAccount) (account *models.Account, err error)
	DescribeAccounts(ctx context.Context, params *params.DescribeAccounts) (pagination *database.Pagination[*models.Account], err error)
	DescribeAccount(ctx context.Context, params *params.DescribeAccount) (account *models.Account, err error)
	AuthenticateAccount(ctx context.Context, params *params.AuthenticateAccount) (account *models.Account, err error)
	UpdateAccount(ctx context.Context, params *params.UpdateAccount) (account *models.Account, err error)

	UpdatePassword(ctx context.Context, params *params.UpdatePassword) (err error)
	ResetPassword(ctx context.Context, params *params.ResetPassword) (err error)
}
