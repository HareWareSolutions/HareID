package controllers

import (
	"HareID/internal/middleware"
	"HareID/internal/responses"
	"HareID/internal/services"
	"errors"
	"net/http"
	"strconv"
)

type NotificationsController struct {
	services services.Services
}

// GetAll retrieves all notifications for a user
// @Summary      Get user notifications
// @Description  Retrieve a list of all notifications for a specific user
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id  path      int  true  "User ID"
// @Success      200      {array}   models.Notification
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/{user_id}/notifications [get]
func (c *NotificationsController) GetAll(w http.ResponseWriter, r *http.Request) {

	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("userkey not found in context"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	userID, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if requestUserID != userID {
		responses.Error(w, http.StatusUnauthorized, errors.New("only the notifications owner can se the notifications"))
		return
	}

	notifications, err := c.services.Notifications.GetAll(r.Context(), requestUserID, userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, notifications)
}

// GetByID retrieves a specific notification
// @Summary      Get notification by ID
// @Description  Retrieve details of a specific notification
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id          path      int  true  "User ID"
// @Param        notification_id  path      int  true  "Notification ID"
// @Success      200              {object}  models.Notification
// @Failure      400              {object}  map[string]string
// @Failure      401              {object}  map[string]string
// @Failure      500              {object}  map[string]string
// @Router       /users/{user_id}/notifications/{notification_id} [get]
func (c *NotificationsController) GetByID(w http.ResponseWriter, r *http.Request) {
	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("userkey not found in context"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	userID, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if requestUserID != userID {
		responses.Error(w, http.StatusUnauthorized, errors.New("only the notifications owner can se the notifications"))
		return
	}

	notificationID, err := strconv.ParseUint(r.PathValue("notification_id"), 10, 64)

	notifications, err := c.services.Notifications.GetByID(r.Context(), requestUserID, userID, notificationID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, notifications)
}

// Delete removes a notification
// @Summary      Delete notification
// @Description  Remove a notification from the system
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id          path      int  true  "User ID"
// @Param        notification_id  path      int  true  "Notification ID"
// @Success      200              {object}  map[string]int
// @Failure      400              {object}  map[string]string
// @Failure      401              {object}  map[string]string
// @Failure      500              {object}  map[string]string
// @Router       /users/{user_id}/notifications/{notification_id} [delete]
func (c *NotificationsController) Delete(w http.ResponseWriter, r *http.Request) {
	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("userkey not found in context"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	userID, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	if requestUserID != userID {
		responses.Error(w, http.StatusUnauthorized, errors.New("only the notifications owner can se the notifications"))
		return
	}

	notificationID, err := strconv.ParseUint(r.PathValue("notification_id"), 10, 64)

	notifications, err := c.services.Notifications.Delete(r.Context(), requestUserID, userID, notificationID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, notifications)
}
