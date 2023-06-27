package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alvinatthariq/walletsvc/entity"
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

func (c *controller) EnableWallet(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")

	token, err := getTokenFromAuth(authorization)
	if err != nil {
		httpRespError(w, r, err, http.StatusUnauthorized)
		return
	}

	wallet, err := c.domain.EnableWallet(token)
	if err != nil {
		if errors.Is(err, entity.ErrorInvalidAuthToken) {
			httpRespError(w, r, err, http.StatusUnauthorized)
			return
		} else if errors.Is(err, entity.ErrorWalletAlreadyEnabled) {
			httpRespError(w, r, err, http.StatusBadRequest)
			return
		}

		httpRespError(w, r, err, http.StatusInternalServerError)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, wallet, nil)
}

func getTokenFromAuth(auth string) (token string, err error) {
	authSplit := strings.Split(auth, " ")
	if len(authSplit) < 2 {
		return token, entity.ErrorInvalidAuthToken
	}

	token = authSplit[1]

	return token, nil
}

func (c *controller) GetWallet(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")

	token, err := getTokenFromAuth(authorization)
	if err != nil {
		httpRespError(w, r, err, http.StatusUnauthorized)
		return
	}

	wallet, err := c.domain.GetWallet(token)
	if err != nil {
		if errors.Is(err, entity.ErrorWalletNotFound) {
			httpRespError(w, r, err, http.StatusNotFound)
			return
		}

		httpRespError(w, r, err, http.StatusInternalServerError)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, wallet, nil)
}
