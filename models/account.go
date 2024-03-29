package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/fox-gonic/fox/database"
	"gorm.io/gorm"
)

// UserStatus user status
type UserStatus string

// UserStatus enum
const (
	UserPendingStatus   UserStatus = "pending"
	UserActivatedStatus UserStatus = "activated"
	UserDisabledStatus  UserStatus = "disabled"
)

// IsValid 是否有效
func (status UserStatus) IsValid() bool {
	switch status {
	case UserPendingStatus, UserActivatedStatus, UserDisabledStatus:
		return true
	}
	return false
}

// Account model
type Account struct {
	database.Model
	Login    string     `json:"login"    gorm:"uniqueIndex;size:255"`
	Email    string     `json:"email"    gorm:"uniqueIndex;size:255"`
	Name     string     `json:"name"     gorm:"size:255"`
	Bio      string     `json:"bio"      gorm:"size:255"`
	Location string     `json:"location" gorm:"size:255"`
	Status   UserStatus `json:"status"   gorm:"size:32"`

	Authentication *Authentication `json:"-"`
	Avatar         string          `json:"avatar" gorm:"-"`
}

// TableName user model table name
func (Account) TableName() string {
	return "accounts"
}

// AfterFind gorm after find callback
func (account *Account) AfterFind(tx *gorm.DB) (err error) {

	// TODO(m) using system configuration
	if len(account.Email) > 0 {
		h := md5.New()
		h.Write([]byte(account.Email))
		account.Avatar = fmt.Sprintf("https://www.gravatar.com/avatar/%s", hex.EncodeToString(h.Sum(nil)))
	}

	return
}

// BeforeDelete gorm before delete callback
func (account *Account) BeforeDelete(tx *gorm.DB) (err error) {

	suffix := fmt.Sprintf(" [deleted-%s]", time.Now().Format("2006-01-02 15:04:05"))

	err = tx.Model(account).Where("`id` = ?", account.ID).UpdateColumn("login", gorm.Expr("concat(login, ?)", suffix)).Error

	return
}
