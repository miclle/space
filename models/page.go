package models

import (
	"database/sql"
	"encoding/json"
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

// PageQuery page unique keys
type PageQuery struct {
	Lang    string
	Version string
}

// Page page meta model
type Page struct {
	ID            int64         `json:"id"             nestedset:"id"             gorm:"primaryKey;autoIncrement"`
	ParentID      sql.NullInt64 `json:"parent_id"      nestedset:"parent_id"      gorm:"index"`
	Lft           int           `json:"-"              nestedset:"lft"`
	Rgt           int           `json:"-"              nestedset:"rgt"`
	Depth         int           `json:"-"              nestedset:"depth"`
	ChildrenCount int           `json:"children_count" nestedset:"children_count"`
	SpaceID       int64         `json:"-"              nestedset:"scope"          gorm:"index"`

	Space           *Space       `json:"space,omitempty"`
	Content         *PageContent `json:"-"`
	FallbackContent *PageContent `json:"-"`

	Children []*Page `json:"children,omitempty" gorm:"-"`
	Parents  []*Page `json:"parents,omitempty"  gorm:"-"`
}

// TableName user model table name
func (Page) TableName() string {
	return "space_pages"
}

// AfterFind gorm after find callback
func (page *Page) AfterFind(tx *gorm.DB) error {
	if page.Content == nil {
		page.Content = page.FallbackContent
	}

	if v, exist := tx.InstanceGet("query"); exist {
		if query, ok := v.(*PageQuery); ok {
			where := tx.Omit("body", "html")
			if query.Lang != "" {
				where = where.Where(&PageContent{Lang: query.Lang})
			}
			if query.Version != "" {
				where = where.Where(&PageContent{Version: query.Version})
			}

			err := tx.InstanceSet("query", nil).
				Joins("Content", where).
				Where("`lft` < ? AND `rgt` > ?", page.Lft, page.Rgt).
				Order("`lft` ASC").
				Find(&page.Parents).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// MarshalJSON implement
func (page *Page) MarshalJSON() ([]byte, error) {
	type Alias Page
	return json.Marshal(&struct {
		*Alias
		ParentID int64 `json:"parent_id"`
		*PageContent
		Space *Space `json:"space,omitempty"`
	}{
		ParentID:    page.ParentID.Int64,
		Alias:       (*Alias)(page),
		PageContent: (*PageContent)(page.Content),
		Space:       page.Space,
	})
}

// PageContent page version model
type PageContent struct {
	ID int64 `json:"-" gorm:"primaryKey"`

	SpaceID   int64 `json:"-" gorm:"index"`
	CreatorID int64 `json:"-" gorm:"index"`

	PageID  int64  `json:"-"       gorm:"uniqueIndex:page_content"`
	Lang    string `json:"lang"    gorm:"uniqueIndex:page_content;size:32"`
	Version string `json:"version" gorm:"uniqueIndex:page_content;size:64"`

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
func (PageContent) TableName() string {
	return "space_page_contents"
}

// Pages pages type
type Pages []*Page

// Build tree
func (nodes *Pages) Build() Pages {

	var (
		root    = Pages{}
		nodeMap = make(map[int64]*Page)
	)

	for _, node := range *nodes {
		nodeMap[node.ID] = node
		if node.ParentID.Int64 == 0 {
			root = append(root, node)
		}
	}

	for _, node := range *nodes {
		if parent, exists := nodeMap[node.ParentID.Int64]; exists {
			parent.Children = append(parent.Children, node)
		}
	}

	return root
}
