package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/alvinatthariq/walletsvc/entity"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (d *domain) InitAccountWallet(customerID string) (customerToken entity.CustomerToken, err error) {
	tx := d.gorm.First(&customerToken, "customer_id = ?", customerID)
	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// generate token if not exist
			customerToken = entity.CustomerToken{
				CustomerID: customerID,
				Token:      generateSecureToken(42),
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			}

			// create to db
			tx := d.gorm.Create(&customerToken)
			err = tx.Error
			if err != nil {
				var mysqlError *mysql.MySQLError
				if errors.As(err, &mysqlError) {
					// check duplicate constraint
					if mysqlError.Number == entity.CodeMySQLDuplicateEntry {
						return customerToken, entity.ErrorAccountAlreadyExist
					}
				}

				return customerToken, err
			}

			return customerToken, nil
		}

		return customerToken, err
	}

	return customerToken, nil
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
