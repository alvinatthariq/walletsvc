package controllers

import (
	"errors"
	"net/http"
	"strconv"
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

	token := strings.Replace(authorization, "Token ", "", -1)

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

func (c *controller) GetWallet(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")

	token := strings.Replace(authorization, "Token ", "", -1)

	wallet, err := c.domain.GetWallet(token)
	if err != nil {
		if errors.Is(err, entity.ErrorInvalidAuthToken) {
			httpRespError(w, r, err, http.StatusUnauthorized)
			return
		} else if errors.Is(err, entity.ErrorWalletNotFound) {
			httpRespError(w, r, err, http.StatusNotFound)
			return
		}

		httpRespError(w, r, err, http.StatusInternalServerError)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, wallet, nil)
}

func (c *controller) GetWalletTransaction(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")

	token := strings.Replace(authorization, "Token ", "", -1)

	transactions, err := c.domain.GetWalletTransaction(token)
	if err != nil {
		if errors.Is(err, entity.ErrorInvalidAuthToken) {
			httpRespError(w, r, err, http.StatusUnauthorized)
			return
		} else if errors.Is(err, entity.ErrorWalletNotFound) {
			httpRespError(w, r, err, http.StatusNotFound)
			return
		} else if errors.Is(err, entity.ErrorWalletDisabled) {
			httpRespError(w, r, err, http.StatusNotFound)
			return
		}

		httpRespError(w, r, err, http.StatusInternalServerError)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, transactions, nil)
}

func (c *controller) DisableWallet(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")

	token := strings.Replace(authorization, "Token ", "", -1)

	wallet, err := c.domain.DisableWallet(token)
	if err != nil {
		if errors.Is(err, entity.ErrorInvalidAuthToken) {
			httpRespError(w, r, err, http.StatusUnauthorized)
			return
		} else if errors.Is(err, entity.ErrorWalletNotFound) {
			httpRespError(w, r, err, http.StatusBadRequest)
			return
		} else if errors.Is(err, entity.ErrorWalletDisabled) {
			httpRespError(w, r, err, http.StatusBadRequest)
			return
		}

		httpRespError(w, r, err, http.StatusInternalServerError)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, wallet, nil)
}

func (c *controller) CreateWalletDeposit(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")

	token := strings.Replace(authorization, "Token ", "", -1)

	amountStr := r.FormValue("amount")
	refID := r.FormValue("reference_id")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		httpRespError(w, r, err, http.StatusBadRequest)
		return
	}

	deposit, err := c.domain.CreateWalletDeposit(token, amount, refID)
	if err != nil {
		if errors.Is(err, entity.ErrorInvalidAuthToken) {
			httpRespError(w, r, err, http.StatusUnauthorized)
			return
		} else if errors.Is(err, entity.ErrorWalletNotFound) {
			httpRespError(w, r, err, http.StatusBadRequest)
			return
		} else if errors.Is(err, entity.ErrorWalletDisabled) {
			httpRespError(w, r, err, http.StatusBadRequest)
			return
		} else if errors.Is(err, entity.ErrorDepositReferenceIDMustBeUnique) {
			httpRespError(w, r, err, http.StatusBadRequest)
			return
		}

		httpRespError(w, r, err, http.StatusInternalServerError)
		return
	}

	httpRespSuccess(w, r, http.StatusOK, deposit, nil)
}
