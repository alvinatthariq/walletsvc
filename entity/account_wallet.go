package entity

import (
	"time"
)

type AccountWallet struct {
	CustomerID string    `json:"customer_id" gorm:"primaryKey;type:varchar(36)"`
	Token      string    `json:"token" gorm:"type:varchar(100)"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:datetime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:datetime"`
}
