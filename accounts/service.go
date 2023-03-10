package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/fox-gonic/fox/database"
	"github.com/miclle/space/accounts/params"
	"github.com/miclle/space/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

// NewService return default implement spaces service
func NewService(database *database.Database) (Service, error) {

	service := &service{
		Database: database,
	}

	return service, nil
}

var _ Service = &service{}

type service struct {
	Database *database.Database
}

func (s *service) CreateAccount(ctx context.Context, params *params.CreateAccount) (*models.Account, error) {
	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err := database.Transaction(func(tx *gorm.DB) error {

		account = &models.Account{
			Login:    params.Login,
			Name:     params.Name,
			Email:    params.Email,
			Bio:      params.Bio,
			Location: params.Location,
		}
		if err := tx.Create(account).Error; err != nil {
			return err
		}

		authentication := &models.Authentication{
			AccountID: account.ID,
			Password:  params.Password,
		}
		if err := tx.Create(authentication).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *service) DescribeAccounts(ctx context.Context, params *params.DescribeAccounts) (*database.Pagination[*models.Account], error) {
	var (
		database   = s.Database.WithContext(ctx)
		pagination = &params.Pagination
	)

	if q := strings.TrimSpace(params.Q); q != "" {
		like := fmt.Sprintf("%%%s%%", q)
		database = database.Where("`accounts`.`name` LIKE ? OR `accounts`.`email` LIKE ?", like, like)
	}

	if err := database.Model(&pagination.Items).Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	// Pagination
	database = database.Scopes(pagination.Paginate())

	if err := database.Find(&pagination.Items).Error; err != nil {
		return nil, err
	}

	return pagination, nil
}

func (s *service) DescribeAccount(ctx context.Context, params *params.DescribeAccount) (*models.Account, error) {
	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	if params.ID == 0 && len(params.Login) == 0 {
		return nil, errors.New("describe account params invalid")
	}

	if params.ID > 0 {
		database = database.Where("`id` = ?", params.ID)
	}

	if len(params.Login) > 0 {
		database = database.Where("`login` = ?", params.Login)
	}

	err := database.First(&account).Error
	if err != nil {
		return nil, err
	}

	return account, err
}

func (s *service) AuthenticateAccount(ctx context.Context, params *params.AuthenticateAccount) (*models.Account, error) {

	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err := database.Preload("Authentication").Where("`login` = ?", params.Login).First(&account).Error
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(account.Authentication.EncryptedPassword, []byte(params.Password)); err != nil {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	return nil, nil
}

func (s *service) UpdateAccount(ctx context.Context, params *params.UpdateAccount) (*models.Account, error) {
	return nil, nil
}

func (s *service) UpdatePassword(ctx context.Context, params *params.UpdatePassword) error {
	return nil
}

func (s *service) ResetPassword(ctx context.Context, params *params.ResetPassword) error {
	return nil
}
