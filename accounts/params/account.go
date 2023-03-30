package params

import (
	"errors"

	"github.com/fox-gonic/fox/database"
	"github.com/miclle/space/models"
)

var (
	// ErrDescribeAccountParamsInvalid describe account params invalid
	ErrDescribeAccountParamsInvalid = errors.New("describe account params invalid")
)

// CreateAccount create account params
type CreateAccount struct {
	Login    string `binding:"required" validate:"required"`
	Email    string `binding:"required" validate:"required,email"`
	Password string `binding:"required" validate:"required,min=8,max=56"`
	Name     string
	Bio      string
	Location string
	Status   models.UserStatus
}

// DescribeAccounts describe accounts params
type DescribeAccounts struct {
	database.Pagination[*models.Account]
	Q string
}

// DescribeAccount describe account params
type DescribeAccount struct {
	ID    int64
	Login string
	Email string
}

// IsValid return describe account params is valid
func (params *DescribeAccount) IsValid() error {
	if params.ID == 0 && len(params.Login) == 0 && len(params.Email) == 0 {
		return ErrDescribeAccountParamsInvalid
	}
	return nil
}

// AuthenticateAccount authenticate account params
type AuthenticateAccount struct {
	Login    string
	Email    string
	Password string
	ClientIP string
}

// UpdateAccount update account params
type UpdateAccount struct {
	Login    string
	Name     *string
	Bio      *string
	Location *string
}

// CreateUnlockToken create unlock account token params
type CreateUnlockToken struct {
	Email string
}

// Unlock account token
type Unlock struct {
	Login string
	Token string
}

// UpdatePassword update password params
type UpdatePassword struct {
	ID          int64
	Login       string
	Email       string
	Password    string
	NewPassword string
}

// CreateResetPasswordToken create reset account password token params
type CreateResetPasswordToken struct {
	Email string
}

// ResetPassword reset password params
type ResetPassword struct {
	Token    string
	Password string
}
