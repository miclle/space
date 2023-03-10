package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Authentication model
type Authentication struct {
	ID                int64  `gorm:"primaryKey"`
	AccountID         int64  `gorm:"uniqueIndex"`                            // Account ID
	Password          string `json:"-"                  gorm:"-"`            // Password
	EncryptedPassword []byte `json:"-"                  gorm:"default:NULL"` // Encrypted password
	CurrentSignInAt   int64  `json:"current_sign_in_at"`                     // Current sign in time
	CurrentSignInIP   string `json:"current_sign_in_ip" gorm:"size:255"`     // Current sign in ip
	LastSignInAt      int64  `json:"last_sign_in_at"`                        // Last sign in time
	LastSignInIP      string `json:"last_sign_in_ip"    gorm:"size:255"`     // Last sign in ip
	FailedAttempts    int    `json:"failed_attempts"`                        // Failed attempt count
	UnlockToken       []byte `json:"unlock_token"       gorm:"default:NULL"` // Lock token
	LockedAt          int64  `json:"locked_at"`                              // Lock at time
	ResetToken        []byte `json:"reset_token"        gorm:"default:NULL"` // Reset token
	ResetAt           int64  `json:"reset_at"`                               // Reset at time
}

// TableName user model table name
func (Authentication) TableName() string {
	return "authentications"
}

// BeforeCreate gorm before create callback
func (auth *Authentication) BeforeCreate(tx *gorm.DB) (err error) {
	auth.EncryptedPassword, err = bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
	return
}
