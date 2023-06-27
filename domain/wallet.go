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
	err = d.gorm.First(&accountWallet, "customer_id = ?", customerID).Error
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
			err = d.gorm.Create(&accountWallet).Error
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

	wallet, err = d.getWalletByCustomerID(customerToken.CustomerID)
	if err != nil {
		if errors.Is(err, entity.ErrorWalletNotFound) {
			// create wallet
			wallet = entity.Wallet{
				ID:        uuid.New().String(),
				OwnedBy:   customerToken.CustomerID,
				Status:    "enabled",
				EnabledAt: time.Now().UTC(),
				Balance:   0,
			}

			// create to db
			err = d.gorm.Create(&wallet).Error
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
		err = d.gorm.Save(wallet).Error
		if err != nil {
			return wallet, err
		}
	}

	return wallet, nil
}

func (d *domain) DisableWallet(token string) (wallet entity.Wallet, err error) {
	customerToken, err := d.getAccountWalletByToken(token)
	if err != nil {
		return wallet, err
	}

	wallet, err = d.getWalletByCustomerID(customerToken.CustomerID)
	if err != nil {
		return wallet, err
	}

	// update status to disabled
	wallet.Status = "disabled"
	err = d.gorm.Save(wallet).Error
	if err != nil {
		return wallet, err
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
	err = d.gorm.First(&accountWallet, "token = ?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return accountWallet, entity.ErrorInvalidAuthToken
		}

		return accountWallet, err
	}

	return accountWallet, nil
}

func (d *domain) getWalletByCustomerID(customerID string) (wallet entity.Wallet, err error) {
	err = d.gorm.First(&wallet, "owned_by = ?", customerID).Error
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

func (d *domain) CreateWalletDeposit(token string, amount float64, refID string) (deposit entity.Deposit, err error) {
	customerToken, err := d.getAccountWalletByToken(token)
	if err != nil {
		return deposit, err
	}

	wallet, err := d.getWalletByCustomerID(customerToken.CustomerID)
	if err != nil {
		return deposit, err
	}

	if wallet.Status == "disabled" {
		return deposit, entity.ErrorWalletDisabled
	}

	deposit = entity.Deposit{
		ID:          uuid.New().String(),
		DepositedBy: customerToken.CustomerID,
		Status:      "success",
		DepositedAt: time.Now().UTC(),
		Amount:      amount,
		ReferenceID: refID,
	}

	if deposit.Amount <= 0 {
		return deposit, entity.ErrorDepositAmountMustBeGreaterThan0
	}

	err = d.gorm.Transaction(func(tx *gorm.DB) error {
		// create transaction
		transaction := entity.Transaction{
			ID:              deposit.ID,
			ReferenceID:     deposit.ReferenceID,
			Type:            "deposit",
			OwnedBy:         deposit.DepositedBy,
			Amount:          deposit.Amount,
			PreviousBalance: wallet.Balance,
			CurrentBalance:  wallet.Balance + amount,
			CreatedAt:       deposit.DepositedAt,
		}
		if err := tx.Create(transaction).Error; err != nil {
			var mysqlError *mysql.MySQLError
			if errors.As(err, &mysqlError) {
				// check duplicate constraint
				if mysqlError.Number == entity.CodeMySQLDuplicateEntry {
					return entity.ErrorDepositReferenceIDMustBeUnique
				}
			}

			return err
		}

		// update wallet balance
		wallet.Balance += amount
		if err := tx.Save(wallet).Error; err != nil {
			return err
		}

		return nil
	})

	return deposit, err
}

func (d *domain) CreateWalletWithdraw(token string, amount float64, refID string) (withdraw entity.Withdraw, err error) {
	customerToken, err := d.getAccountWalletByToken(token)
	if err != nil {
		return withdraw, err
	}

	wallet, err := d.getWalletByCustomerID(customerToken.CustomerID)
	if err != nil {
		return withdraw, err
	}

	if wallet.Status == "disabled" {
		return withdraw, entity.ErrorWalletDisabled
	}

	withdraw = entity.Withdraw{
		ID:          uuid.New().String(),
		WithdrawnBy: customerToken.CustomerID,
		Status:      "success",
		WithdrawnAt: time.Now().UTC(),
		Amount:      amount,
		ReferenceID: refID,
	}

	if withdraw.Amount > wallet.Balance {
		return withdraw, entity.ErrorBalanceNotEnough
	}

	err = d.gorm.Transaction(func(tx *gorm.DB) error {
		// create transaction
		transaction := entity.Transaction{
			ID:              withdraw.ID,
			ReferenceID:     withdraw.ReferenceID,
			Type:            "withdraw",
			OwnedBy:         withdraw.WithdrawnBy,
			Amount:          withdraw.Amount,
			PreviousBalance: wallet.Balance,
			CurrentBalance:  wallet.Balance - amount,
			CreatedAt:       withdraw.WithdrawnAt,
		}
		if err := tx.Create(transaction).Error; err != nil {
			var mysqlError *mysql.MySQLError
			if errors.As(err, &mysqlError) {
				// check duplicate constraint
				if mysqlError.Number == entity.CodeMySQLDuplicateEntry {
					return entity.ErrorWithdrawReferenceIDMustBeUnique
				}
			}

			return err
		}

		// update wallet balance
		wallet.Balance -= amount
		if err := tx.Save(wallet).Error; err != nil {
			return err
		}

		return nil
	})

	return withdraw, err
}
