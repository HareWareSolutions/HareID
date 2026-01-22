package controllers

import (
	"HareID/internal/models"
	"HareID/internal/responses"
	"HareID/internal/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type UsersController struct {
	services services.Services
}

func (c *UsersController) Create(w http.ResponseWriter, r *http.Request) {

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	newUser, err := c.services.Users.Create(r.Context(), user)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, newUser)
}

func (c *UsersController) GetAll(w http.ResponseWriter, r *http.Request) {

	users, err := c.services.Users.GetAll(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

func (c *UsersController) GetByID(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	user, err := c.services.Users.GetByID(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func (c *UsersController) GetUserTeam(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	teamMember, err := c.services.TeamMembers.GetByUserID(r.Context(), userID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, teamMember)
}

func (c *UsersController) Update(w http.ResponseWriter, r *http.Request) {

	requestUserIDString, ok := r.Context().Value("user_id").(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("UserKey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	affectedRows, err := c.services.Users.Update(r.Context(), userID, requestUserID, user)
	if err != nil {
		responses.Error(w, http.StatusForbidden, err)
		return
	}

	data := map[string]uint64{
		"affected_rows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}

func (c *UsersController) Delete(w http.ResponseWriter, r *http.Request) {

	requestUserIDString, ok := r.Context().Value("user_id").(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("UserKey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	affectedRows, err := c.services.Users.Delete(r.Context(), userID, requestUserID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]uint64{
		"affectedRows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}
