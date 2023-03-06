package models

import "gorm.io/plugin/soft_delete"

// Version model
type Version struct {
	ID      int64  `json:"id"       gorm:"primaryKey"`
	SpaceID int64  `json:"space_id" gorm:"uniqueIndex:space_version"`
	Name    string `json:"name"     gorm:"uniqueIndex:space_version"`

	CreatedAt int64                 `json:"created_at"`
	UpdatedAt int64                 `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
