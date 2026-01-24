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

// Create creates a new subscription
// @Summary      Create subscription
// @Description  Create a new subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        subscription  body      models.Subscription  true  "Subscription Data"
// @Success      201           {object}  models.Subscription
// @Failure      500           {object}  map[string]string
// @Router       /subscriptions [post]
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

// GetAll retrieves all subscriptions
// @Summary      Get all subscriptions
// @Description  Retrieve a list of all subscriptions
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200   {array}   models.Subscription
// @Failure      500   {object}  map[string]string
// @Router       /subscriptions [get]
func (c *SubscriptionsController) GetAll(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := c.services.Subscriptions.GetAll(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, subscriptions)
}

// GetBySubscriptionID retrieves a subscription by ID
// @Summary      Get subscription by ID
// @Description  Retrieve details of a specific subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        subscription_id  path      string  true  "Subscription ID"
// @Success      200              {object}  models.Subscription
// @Failure      400              {object}  map[string]string
// @Failure      500              {object}  map[string]string
// @Router       /subscriptions/{subscription_id} [get]
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

// Update modifies an existing subscription
// @Summary      Update subscription
// @Description  Update details of an existing subscription
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        subscription_id  path      string               true  "Subscription ID"
// @Param        subscription     body      models.Subscription  true  "Subscription Update Data"
// @Success      200              {object}  map[string]uint64
// @Failure      400              {object}  map[string]string
// @Failure      403              {object}  map[string]string
// @Failure      500              {object}  map[string]string
// @Router       /subscriptions/{subscription_id} [patch]
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

// Delete removes a subscription
// @Summary      Delete subscription
// @Description  Remove a subscription from the system
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        subscription_id  path      string  true  "Subscription ID"
// @Success      200              {object}  map[string]uint64
// @Failure      400              {object}  map[string]string
// @Failure      403              {object}  map[string]string
// @Failure      500              {object}  map[string]string
// @Router       /subscriptions/{subscription_id} [delete]
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
