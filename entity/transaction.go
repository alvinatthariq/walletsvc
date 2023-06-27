package entity

import (
	"time"
)

type Transaction struct {
	ID              string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	ReferenceID     string    `json:"reference_id" gorm:"type:varchar(36);index:refid_type_idx,unique"`
	Type            string    `json:"type" gorm:"type:varchar(36);index:refid_type_idx,unique"`
	OwnedBy         string    `json:"owned_by" gorm:"type:varchar(36)"`
	Amount          float64   `json:"amount" gorm:"type:decimal(19,2)"`
	PreviousBalance float64   `json:"previous_balance" gorm:"type:decimal(19,2)"`
	CurrentBalance  float64   `json:"current_balance" gorm:"type:decimal(19,2)"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:datetime"`
}

type Deposit struct {
	ID          string
	DepositedBy string
	Status      string
	DepositedAt time.Time
	Amount      float64
	ReferenceID string
}
