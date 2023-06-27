package entity

import (
	"time"
)

type Wallet struct {
	ID        string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	OwnedBy   string    `json:"owned_by" gorm:"type:varchar(36)"`
	Status    string    `json:"status" gorm:"type:varchar(50)"`
	EnabledAt time.Time `json:"enabled_at" gorm:"type:datetime"`
	Balance   float64   `json:"balance" gorm:"type:decimal(19,2)"`
}
