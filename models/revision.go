package models

// Revision model
type Revision struct {
	ID        int64 `json:"id"         gorm:"primaryKey"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`

	OwnerID   int64  `json:"owner_id"`
	OwnerType string `json:"owner_type"`
}
