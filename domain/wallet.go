package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/alvinatthariq/walletsvc/entity"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

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

func (d *domain) EnableWallet(token string) (wallet entity.Wallet, err error) {
	var customerToken entity.CustomerToken
	tx := d.gorm.First(&customerToken, "token = ?", token)
	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return wallet, entity.ErrorInvalidAuthToken
		}

		return wallet, err
	}

	tx = d.gorm.First(&wallet, "owned_by = ?", customerToken.CustomerID)
	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create wallet
			wallet = entity.Wallet{
				ID:        uuid.New().String(),
				OwnedBy:   customerToken.CustomerID,
				Status:    "enabled",
				EnabledAt: time.Now().UTC(),
				Balance:   0,
			}

			// create to db
			tx := d.gorm.Create(&wallet)
			err = tx.Error
			if err != nil {
				var mysqlError *mysql.MySQLError
				if errors.As(err, &mysqlError) {
					// check duplicate constraint
					if mysqlError.Number == entity.CodeMySQLDuplicateEntry {
						return wallet, entity.ErrorWalletAlreadyExist
					}
				}

				return wallet, err
			}

			return wallet, nil
		}

		return wallet, err
	}

	if wallet.Status == "enabled" {
		// return error if already enabled
		return wallet, entity.ErrorWalletAlreadyEnabled
	} else {
		// update wallet status to enabled
		wallet.Status = "enabled"
		wallet.EnabledAt = time.Now().UTC()
		tx = d.gorm.Save(wallet)
		err = tx.Error
		if err != nil {
			return wallet, err
		}
	}

	return wallet, nil
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
