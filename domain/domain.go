package domain

import (
	"github.com/alvinatthariq/walletsvc/entity"
	"github.com/go-redis/redis"

	"gorm.io/gorm"
)

type domain struct {
	gorm        *gorm.DB
	redisClient *redis.Client
}

type DomainItf interface {
	InitAccountWallet(customerID string) (accountWallet entity.AccountWallet, err error)
	EnableWallet(token string) (wallet entity.Wallet, err error)
	DisableWallet(token string) (wallet entity.Wallet, err error)
	GetWallet(token string) (wallet entity.Wallet, err error)
	GetWalletTransaction(token string) (transactions []entity.Transaction, err error)
	CreateWalletDeposit(token string, amount float64, refID string) (deposit entity.Deposit, err error)
	CreateWalletWithdraw(token string, amount float64, refID string) (withdraw entity.Withdraw, err error)
}

func Init(gorm *gorm.DB, redisClient *redis.Client) DomainItf {
	return &domain{
		gorm:        gorm,
		redisClient: redisClient,
	}
}
