package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/fox-gonic/fox/database"
	"gorm.io/gorm"
)

var (
	// ErrSpaceStatusIsInvalid space status is invalid
	ErrSpaceStatusIsInvalid = errors.New("space status is invalid")
)

// SpaceStatus space status
type SpaceStatus string

// SpaceStatus enum
const (
	SpaceStatusOffline SpaceStatus = "offline"
	SpaceStatusOnline  SpaceStatus = "online"
)

// IsValid return space status is valid
func (t SpaceStatus) IsValid() error {
	switch t {
	case
		SpaceStatusOffline, SpaceStatusOnline:
		return nil
	default:
		return ErrSpaceStatusIsInvalid
	}
}

// Space model
type Space struct {
	database.Model
	Name         string      `json:"name"          gorm:"uniqueIndex;size:128"`
	Key          string      `json:"key"           gorm:"uniqueIndex;size:128"`
	Lang         string      `json:"lang"          gorm:"size:32"` // default lang // TODO(m) enum type
	FallbackLang string      `json:"fallback_lang" gorm:"size:32"` // fallback lang
	HomepageID   int64       `json:"homepage_id"   gorm:"index"`
	Description  string      `json:"description"`
	Avatar       string      `json:"avatar"`
	Status       SpaceStatus `json:"status"        gorm:"index;size:32"`
	CreatorID    int64       `json:"-"`

	Homepage *Page `json:"homepage,omitempty" gorm:"foreignKey:HomepageID;references:PageID"`
}

// TableName user model table name
func (Space) TableName() string {
	return "spaces"
}

// BeforeDelete gorm before delete callback
func (space *Space) BeforeDelete(tx *gorm.DB) (err error) {
	now := time.Now().Unix()
	updates := map[string]interface{}{
		"`name`": gorm.Expr("concat(`name`, ?)", fmt.Sprintf(" [deleted-%d]", now)),
		"`key`":  gorm.Expr("concat(`key`, ?)", fmt.Sprintf(" [deleted-%d]", now)),
	}
	err = tx.Model(space).Where("`id` = ?", space.ID).UpdateColumns(updates).Error
	return
}
