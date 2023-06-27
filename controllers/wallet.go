package controllers

import (
	"net/http"
)

func (c *controller) InitAccountWallet(w http.ResponseWriter, r *http.Request) {
	customerID := r.FormValue("customer_xid")

	customerToken, err := c.domain.InitAccountWallet(customerID)
	if err != nil {
		httpRespError(w, r, err, http.StatusInternalServerError)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, customerToken, nil)
}
