package controllers

import (
	"HareCRM/internal/models"
	"HareCRM/internal/responses"
	"HareCRM/internal/services"
	"encoding/json"
	"net/http"
)

type LoginController struct {
	services services.Services
}

func (c *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	token, err := c.services.Login.Login(r.Context(), user.GoogleSub)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	data := map[string]string{
		"token": token,
	}

	responses.JSON(w, http.StatusOK, data)
}
