package params

import (
	"github.com/fox-gonic/fox/database"
	"github.com/miclle/space/models"
)

// CreateAccount create account params
type CreateAccount struct {
	Login    string `binding:"required" validate:"required"`
	Email    string `binding:"required" validate:"required,email"`
	Password string `binding:"required" validate:"required,min=8,max=56"`
	Name     string
	Bio      string
	Location string
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
}

// AuthenticateAccount authenticate account params
type AuthenticateAccount struct {
	Login    string
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
	Login string
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
	Password    string
	NewPassword string
}

// CreateResetPasswordToken create reset account password token params
type CreateResetPasswordToken struct {
	Login string
}

// ResetPassword reset password params
type ResetPassword struct {
	Login    string
	Token    string
	Password string
}
