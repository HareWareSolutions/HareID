package controllers

import (
	"HareID/internal/authentication"
	"HareID/internal/responses"
	"HareID/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type CheckoutController struct {
	services services.Services
}

type CreateCheckoutRequest struct {
	PriceID    string `json:"price_id"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

func (c *CheckoutController) CreateSession(w http.ResponseWriter, r *http.Request) {
	requestUserID, err := authentication.GetTokenUserID(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	userID, err := strconv.ParseUint(requestUserID, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	var req CreateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if req.PriceID == "" || req.SuccessURL == "" || req.CancelURL == "" {
		responses.Error(w, http.StatusBadRequest, errors.New("price_id, success_url and cancel_url are required"))
		return
	}

	checkoutURL, err := c.services.Checkout.CreateCheckoutSession(r.Context(), userID, req.PriceID, req.SuccessURL, req.CancelURL)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, map[string]string{"url": checkoutURL})
}
