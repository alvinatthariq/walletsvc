package controllers

import (
	"github.com/alvinatthariq/walletsvc/domain"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type controller struct {
	gorm   *gorm.DB
	router *mux.Router
	domain domain.DomainItf
}

func Init(gorm *gorm.DB, router *mux.Router, domain domain.DomainItf) {
	var c *controller

	c = &controller{
		gorm:   gorm,
		router: router,
		domain: domain,
	}

	c.Serve()
}

func (c *controller) Serve() {
	c.router.HandleFunc("/api/v1/init", c.InitAccountWallet).Methods("POST")
	c.router.HandleFunc("/api/v1/wallet", c.EnableWallet).Methods("POST")
	c.router.HandleFunc("/api/v1/wallet", c.DisableWallet).Methods("PATCH")
	c.router.HandleFunc("/api/v1/wallet", c.GetWallet).Methods("GET")
	c.router.HandleFunc("/api/v1/wallet/transactions", c.GetWalletTransaction).Methods("GET")
	c.router.HandleFunc("/api/v1/wallet/deposits", c.CreateWalletDeposit).Methods("POST")
	// c.router.HandleFunc("/api/v1/wallet/withdrawals", c.CreateWalletWithdraw).Methods("POST")

}
