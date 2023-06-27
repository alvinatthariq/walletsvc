package entity

import (
	"time"
)

type Transaction struct {
	ReferenceID     string    `json:"reference_id" gorm:"primaryKey;type:varchar(36)"`
	Type            string    `json:"type" gorm:"type:varchar(36)"`
	OwnedBy         string    `json:"owned_by" gorm:"type:varchar(36)"`
	Amount          float64   `json:"amount" gorm:"type:decimal(19,2)"`
	PreviousBalance float64   `json:"previous_balance" gorm:"type:decimal(19,2)"`
	CurrentBalance  float64   `json:"current_balance" gorm:"type:decimal(19,2)"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime"`
}
