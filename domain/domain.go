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
	GetWallet(token string) (wallet entity.Wallet, err error)
	GetWalletTransaction(token string) (transactions []entity.Transaction, err error)
}

func Init(gorm *gorm.DB, redisClient *redis.Client) DomainItf {
	return &domain{
		gorm:        gorm,
		redisClient: redisClient,
	}
}
