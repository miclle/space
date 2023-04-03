package models

import (
	"database/sql"
	"errors"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

var (
	// ErrPageStatusIsInvalid page status is invalid
	ErrPageStatusIsInvalid = errors.New("page status is invalid")
)

// PageStatus page status
type PageStatus string

// PageStatus enum
const (
	PageStatusDraft      PageStatus = "draft"
	PageStatusPublished  PageStatus = "published"
	PageStatusOffline    PageStatus = "offline"
	PageStatusDeprecated PageStatus = "deprecated"
)

// IsValid return space status is valid
func (t PageStatus) IsValid() error {
	switch t {
	case
		PageStatusDraft, PageStatusPublished, PageStatusOffline, PageStatusDeprecated:
		return nil
	default:
		return ErrPageStatusIsInvalid
	}
}

// PageMeta page meta model
type PageMeta struct {
	ID            int64         `json:"id"             nestedset:"id"             gorm:"primaryKey;autoIncrement"`
	ParentID      sql.NullInt64 `json:"parent_id"      nestedset:"parent_id"      gorm:"index"`
	Rgt           int           `json:"rgt"            nestedset:"rgt"`
	Lft           int           `json:"lft"            nestedset:"lft"`
	Depth         int           `json:"depth"          nestedset:"depth"`
	ChildrenCount int           `json:"children_count" nestedset:"children_count"`

	SpaceID int64  `json:"-"               gorm:"index"`
	Space   *Space `json:"space,omitempty"`
}

// Page page version model
type Page struct {
	ID int64 `json:"-" gorm:"primaryKey"`

	SpaceID      int64 `json:"-" gorm:"index"`
	CreatorID    int64 `json:"-" gorm:"index"`
	ParentPageID int64 `json:"-" gorm:"index"`

	PageID     int64      `json:"id"          gorm:"uniqueIndex:page_content"`
	Lang       string     `json:"lang"        gorm:"uniqueIndex:page_content;size:32"`
	Version    string     `json:"version"     gorm:"uniqueIndex:page_content;size:64"`
	Status     PageStatus `json:"status"      gorm:"size:32"`
	Title      string     `json:"title"       gorm:"size:255;index"`
	ShortTitle string     `json:"short_title" gorm:"size:255;index"`
	Body       string     `json:"body"`
	HTML       string     `json:"html"`

	CreatedAt int64                 `json:"created_at"`
	UpdatedAt int64                 `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Space *Space `json:"space,omitempty"`
}

// TableName user model table name
func (Page) TableName() string {
	return "space_pages"
}

// BeforeCreate gorm before create callback
func (page *Page) BeforeCreate(tx *gorm.DB) (err error) {
	if page.PageID > 0 {
		var p *Page
		err = tx.Where("`space_id` = ? AND `page_id` = ?", page.SpaceID, page.PageID).First(&p).Error
		if err != nil {
			return err
		}
		page.ParentPageID = p.ParentPageID
	}
	return
}

// AfterCreate gorm after create callback
func (page *Page) AfterCreate(tx *gorm.DB) (err error) {
	if page.PageID == 0 {
		tx.Model(page).Update("page_id", page.ID)
	}
	return
}
