package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fox-gonic/fox/database"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/miclle/space/accounts/params"
	"github.com/miclle/space/models"
)

// validate global validator
var validate = validator.New()

// Service for account interface
type Service interface {
	CreateAccount(ctx context.Context, params *params.CreateAccount) (account *models.Account, err error)
	DescribeAccounts(ctx context.Context, params *params.DescribeAccounts) (pagination *database.Pagination[*models.Account], err error)
	DescribeAccount(ctx context.Context, params *params.DescribeAccount) (account *models.Account, err error)
	AuthenticateAccount(ctx context.Context, params *params.AuthenticateAccount) (account *models.Account, err error)
	UpdateAccount(ctx context.Context, params *params.UpdateAccount) (account *models.Account, err error)

	CreateUnlockToken(ctx context.Context, params *params.CreateUnlockToken) (token string, err error)
	Unlock(ctx context.Context, params *params.Unlock) (err error)

	UpdatePassword(ctx context.Context, params *params.UpdatePassword) (err error)

	CreateResetPasswordToken(ctx context.Context, params *params.ResetPassword) (token string, err error)
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

	if err := validate.Struct(params); err != nil {
		return nil, err
	}

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
		database = database.Where("`accounts`.`login` LIKE ? OR `accounts`.`name` LIKE ? OR `accounts`.`email` LIKE ?", like, like, like)
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

	if err := params.IsValid(); err != nil {
		return nil, err
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(http.StatusText(http.StatusUnauthorized))
		}
		return nil, err
	}

	authentication := account.Authentication

	if err := bcrypt.CompareHashAndPassword(authentication.EncryptedPassword, []byte(params.Password)); err != nil {
		authentication.FailedAttempts++
		database.Save(authentication)
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	authentication.FailedAttempts = 0
	authentication.LastSignInAt = authentication.CurrentSignInAt
	authentication.LastSignInIP = authentication.CurrentSignInIP
	authentication.CurrentSignInAt = time.Now().Unix()
	authentication.CurrentSignInIP = params.ClientIP

	database.Save(authentication)

	return account, nil
}

func (s *service) UpdateAccount(ctx context.Context, params *params.UpdateAccount) (*models.Account, error) {
	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err := database.Where("`login` = ?", params.Login).First(&account).Error
	if err != nil {
		return nil, err
	}

	if params.Name != nil {
		account.Name = *params.Name
	}
	if params.Bio != nil {
		account.Bio = *params.Bio
	}
	if params.Location != nil {
		account.Location = *params.Location
	}

	err = database.Save(account).Error

	return account, err
}

func (s *service) CreateUnlockToken(ctx context.Context, params *params.CreateUnlockToken) (token string, err error) {
	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err = database.Preload("Authentication").Where("`login` = ?", params.Login).First(&account).Error
	if err != nil {
		return "", err
	}

	err = database.Transaction(func(tx *gorm.DB) error {
		var (
			password = uuid.New().String()
			now      = jwt.NewNumericDate(time.Now())
			t        = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
				Issuer:    params.Login,
				Subject:   password,
				IssuedAt:  now,
				NotBefore: now,
				ExpiresAt: &jwt.NumericDate{Time: now.Add(time.Hour * 24)},
			})
		)

		encryptedToken, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		err = tx.Model(account.Authentication).UpdateColumn("unlock_token", encryptedToken).Error
		if err != nil {
			return err
		}

		token, err = t.SignedString([]byte("unlock"))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return
}

func (s *service) Unlock(ctx context.Context, params *params.Unlock) (err error) {

	token, err := jwt.Parse(params.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		// if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		// }
		return []byte("unlock"), nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if ok && token.Valid {
		fmt.Printf("claims: %+v\n", claims)
	} else {
		fmt.Println(err)
	}

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, params *params.UpdatePassword) error {
	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err := database.Preload("Authentication").Where("`login` = ?", params.Login).First(&account).Error
	if err != nil {
		return err
	}

	authentication := account.Authentication

	if err := bcrypt.CompareHashAndPassword(authentication.EncryptedPassword, []byte(params.Password)); err != nil {
		return errors.New(http.StatusText(http.StatusUnauthorized))
	}

	authentication.EncryptedPassword, err = bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return database.Save(authentication).Error
}

func (s *service) CreateResetPasswordToken(ctx context.Context, params *params.ResetPassword) (token string, err error) {
	// TODO(m)
	return "", nil
}

func (s *service) ResetPassword(ctx context.Context, params *params.ResetPassword) error {
	// TODO(m)
	return nil
}
