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
