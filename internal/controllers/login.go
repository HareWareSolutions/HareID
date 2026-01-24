package controllers

import (
	"HareID/internal/models"
	"HareID/internal/responses"
	"HareID/internal/services"
	"encoding/json"
	"net/http"
)

type LoginController struct {
	services services.Services
}

// Login authenticates a user
// @Summary      User Login
// @Description  Authenticate user using Google Subject ID and return a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      models.User  true  "User Credentials (only GoogleSub needed)"
// @Success      200          {object}  map[string]string
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /login [post]
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
