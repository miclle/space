package accounts

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fox-gonic/fox/database"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/miclle/space/accounts/params"
	"github.com/miclle/space/models"
)

// validate global validator
var validate = validator.New()

var (
	// SigningMethod jwt signing method
	SigningMethod *jwt.SigningMethodHMAC = jwt.SigningMethodHS256

	// TokenExpires token expires duration
	TokenExpires time.Duration = time.Minute * 15

	// UnlockTokenKey unlock token key
	UnlockTokenKey = "__unlock_token__"

	// ResetPasswordTokenKey reset password token key
	ResetPasswordTokenKey = "__reset_password_token__"
)

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

	CreateResetPasswordToken(ctx context.Context, params *params.CreateResetPasswordToken) (token string, err error)
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

	if len(params.Email) > 0 {
		database = database.Where("`email` = ?", params.Email)
	}

	err := database.Preload("Authentication").First(&account).Error
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

	if len(params.Login) == 0 && len(params.Email) == 0 {
		return nil, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	if len(params.Login) > 0 {
		database = database.Where("`login` = ?", params.Login)
	}

	if len(params.Email) > 0 {
		database = database.Where("`email` = ?", params.Email)
	}

	err := database.Preload("Authentication").First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(http.StatusText(http.StatusUnauthorized))
		}
		return nil, err
	}

	database = s.Database.WithContext(ctx)
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

	if err := database.Save(authentication).Error; err != nil {
		return nil, err
	}

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

func (s *service) CreateUnlockToken(ctx context.Context, params *params.CreateUnlockToken) (string, error) {
	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err := database.Preload("Authentication").Where("`email` = ?", params.Email).First(&account).Error
	if err != nil {
		return "", err
	}

	var (
		password = ulid.Make().String()
		t        = jwt.NewWithClaims(SigningMethod, jwt.RegisteredClaims{
			Issuer:    params.Email,
			ID:        password,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpires)),
		})
	)

	signedStr, err := t.SignedString([]byte(UnlockTokenKey))
	if err != nil {
		return "", err
	}

	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// save unlock token encrypted hash
	err = database.Model(account.Authentication).Update("unlock_token", encryptedToken).Error
	if err != nil {
		return "", err
	}

	var (
		segs  = strings.Split(signedStr, ".")
		token = strings.Join(segs[1:], ".")
	)

	return base64.RawURLEncoding.EncodeToString([]byte(token)), nil
}

func (s *service) Unlock(ctx context.Context, params *params.Unlock) (err error) {

	header := map[string]interface{}{
		"typ": "JWT",
		"alg": SigningMethod.Alg(),
	}

	jsonValue, err := json.Marshal(header)
	if err != nil {
		return err
	}

	tokenBytes, err := base64.RawURLEncoding.DecodeString(params.Token)
	if err != nil {
		return err
	}

	var (
		seg         = jwt.EncodeSegment(jsonValue)
		tokenString = fmt.Sprintf("%s.%s", seg, string(tokenBytes))
	)

	claims := jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(UnlockTokenKey), nil
	})
	if err != nil {
		return
	}

	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err = database.Preload("Authentication").Where("`email` = ?", claims.Issuer).First(&account).Error
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(account.Authentication.UnlockToken, []byte(claims.ID)); err != nil {
		return errors.New(http.StatusText(http.StatusUnauthorized))
	}

	var updates = map[string]interface{}{
		"failed_attempts": 0,
		"unlock_token":    nil,
		"locked_at":       0,
	}
	err = database.Model(account.Authentication).Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *service) UpdatePassword(ctx context.Context, params *params.UpdatePassword) error {
	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	if len(params.NewPassword) == 0 {
		return errors.New("new passowrd is invalid")
	}

	if params.ID > 0 {
		database = database.Where("`id` = ?", params.ID)
	}

	if len(params.Login) > 0 {
		database = database.Where("`login` = ?", params.Login)
	}

	if len(params.Email) > 0 {
		database = database.Where("`email` = ?", params.Email)
	}

	err := database.Preload("Authentication").First(&account).Error
	if err != nil {
		return err
	}

	authentication := account.Authentication
	if err := bcrypt.CompareHashAndPassword(authentication.EncryptedPassword, []byte(params.Password)); err != nil {
		return errors.New(http.StatusText(http.StatusUnauthorized))
	}

	database = s.Database.WithContext(ctx)

	authentication.EncryptedPassword, err = bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return database.Save(authentication).Error
}

func (s *service) CreateResetPasswordToken(ctx context.Context, params *params.CreateResetPasswordToken) (string, error) {
	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err := database.Preload("Authentication").Where("`email` = ?", params.Email).First(&account).Error
	if err != nil {
		return "", err
	}

	var (
		password = ulid.Make().String()
		t        = jwt.NewWithClaims(SigningMethod, jwt.RegisteredClaims{
			Issuer:    params.Email,
			ID:        password,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpires)),
		})
	)

	signedStr, err := t.SignedString([]byte(ResetPasswordTokenKey))
	if err != nil {
		return "", err
	}

	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// save reset token encrypted hash
	err = database.Model(account.Authentication).Update("reset_token", encryptedToken).Error
	if err != nil {
		return "", err
	}

	var (
		segs  = strings.Split(signedStr, ".")
		token = strings.Join(segs[1:], ".")
	)

	return base64.RawURLEncoding.EncodeToString([]byte(token)), nil
}

func (s *service) ResetPassword(ctx context.Context, params *params.ResetPassword) error {
	header := map[string]interface{}{
		"typ": "JWT",
		"alg": SigningMethod.Alg(),
	}

	jsonValue, err := json.Marshal(header)
	if err != nil {
		return err
	}

	tokenBytes, err := base64.RawURLEncoding.DecodeString(params.Token)
	if err != nil {
		return err
	}

	var (
		seg         = jwt.EncodeSegment(jsonValue)
		tokenString = fmt.Sprintf("%s.%s", seg, string(tokenBytes))
	)

	claims := jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(ResetPasswordTokenKey), nil
	})
	if err != nil {
		return err
	}

	var (
		database = s.Database.WithContext(ctx)
		account  *models.Account
	)

	err = database.Preload("Authentication").Where("`email` = ?", claims.Issuer).First(&account).Error
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(account.Authentication.ResetToken, []byte(claims.ID)); err != nil {
		return errors.New(http.StatusText(http.StatusUnauthorized))
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var updates = map[string]interface{}{
		"encrypted_password": encryptedPassword,
	}
	err = database.Model(account.Authentication).Updates(updates).Error
	if err != nil {
		return err
	}

	return nil
}
