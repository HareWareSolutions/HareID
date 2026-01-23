package controllers

import (
	"HareID/internal/models"
	"HareID/internal/responses"
	"HareID/internal/services"
	"encoding/json"
	"errors"
	"net/http"
)

type SubscriptionsController struct {
	services services.Services
}

func (c *SubscriptionsController) Create(w http.ResponseWriter, r *http.Request) {
	var subscription models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	newSubscription, err := c.services.Subscriptions.Create(r.Context(), subscription)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newSubscription)
}

func (c *SubscriptionsController) GetAll(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := c.services.Subscriptions.GetAll(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, subscriptions)
}

func (c *SubscriptionsController) GetBySubscriptionID(w http.ResponseWriter, r *http.Request) {

	subscriptionID := r.PathValue("subscription_id")
	if subscriptionID == "" {
		responses.Error(w, http.StatusBadRequest, errors.New("subscription_id is required"))
		return
	}

	subscription, err := c.services.Subscriptions.GetBySubscriptionID(r.Context(), subscriptionID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, subscription)
}

func (c *SubscriptionsController) Update(w http.ResponseWriter, r *http.Request) {

	subscriptionID := r.PathValue("subscription_id")
	if subscriptionID == "" {
		responses.Error(w, http.StatusBadRequest, errors.New("subscription_id is required"))
		return
	}

	var subscription models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&subscription); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	affectedRows, err := c.services.Subscriptions.Update(r.Context(), subscriptionID, subscription)
	if err != nil {
		responses.Error(w, http.StatusForbidden, err)
		return
	}

	data := map[string]uint64{
		"affected_rows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}

func (c *SubscriptionsController) Delete(w http.ResponseWriter, r *http.Request) {

	subscriptionID := r.PathValue("subscription_id")
	if subscriptionID == "" {
		responses.Error(w, http.StatusBadRequest, errors.New("subscription_id is required"))
		return
	}

	affectedRows, err := c.services.Subscriptions.Delete(r.Context(), subscriptionID)
	if err != nil {
		responses.Error(w, http.StatusForbidden, err)
		return
	}

	data := map[string]uint64{
		"affected_rows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}
