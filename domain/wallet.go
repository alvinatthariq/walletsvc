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

func (d *domain) InitAccountWallet(customerID string) (accountWallet entity.AccountWallet, err error) {
	tx := d.gorm.First(&accountWallet, "customer_id = ?", customerID)
	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// create account wallet if not exist
			accountWallet = entity.AccountWallet{
				CustomerID: customerID,
				Token:      generateSecureToken(42),
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			}

			// create to db
			tx := d.gorm.Create(&accountWallet)
			err = tx.Error
			if err != nil {
				var mysqlError *mysql.MySQLError
				if errors.As(err, &mysqlError) {
					// check duplicate constraint
					if mysqlError.Number == entity.CodeMySQLDuplicateEntry {
						return accountWallet, entity.ErrorAccountWalletAlreadyExist
					}
				}

				return accountWallet, err
			}

			return accountWallet, nil
		}

		return accountWallet, err
	}

	return accountWallet, nil
}

func (d *domain) EnableWallet(token string) (wallet entity.Wallet, err error) {
	customerToken, err := d.getAccountWalletByToken(token)
	if err != nil {
		return wallet, err
	}

	tx := d.gorm.First(&wallet, "owned_by = ?", customerToken.CustomerID)
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

func (d *domain) GetWallet(token string) (wallet entity.Wallet, err error) {
	customerToken, err := d.getAccountWalletByToken(token)
	if err != nil {
		return wallet, err
	}

	wallet, err = d.getWalletByCustomerID(customerToken.CustomerID)
	if err != nil {
		return wallet, err
	}

	return wallet, nil
}

func (d *domain) getAccountWalletByToken(token string) (accountWallet entity.AccountWallet, err error) {
	tx := d.gorm.First(&accountWallet, "token = ?", token)
	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return accountWallet, entity.ErrorInvalidAuthToken
		}

		return accountWallet, err
	}

	return accountWallet, nil
}

func (d *domain) getWalletByCustomerID(customerID string) (wallet entity.Wallet, err error) {
	tx := d.gorm.First(&wallet, "owned_by = ?", customerID)
	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return wallet, entity.ErrorWalletNotFound
		}

		return wallet, err
	}

	return wallet, nil
}

func (d *domain) GetWalletTransaction(token string) (transactions []entity.Transaction, err error) {
	customerToken, err := d.getAccountWalletByToken(token)
	if err != nil {
		return transactions, err
	}

	wallet, err := d.getWalletByCustomerID(customerToken.CustomerID)
	if err != nil {
		return transactions, err
	}

	if wallet.Status == "disabled" {
		return transactions, entity.ErrorWalletDisabled
	}

	tx := d.gorm.Where("owned_by = ?", customerToken.CustomerID).Find(&transactions)
	err = tx.Error

	return transactions, err
}
