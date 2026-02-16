package controllers

import (
	"HareID/internal/middleware"
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

// Create creates a new user
// @Summary      Create a new user
// @Description  Register a new user in the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Registration Data"
// @Success      201   {object}  models.User
// @Failure      500   {object}  map[string]string
// @Router       /users [post]
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

// GetAll retrieves all users
// @Summary      Get all users
// @Description  Retrieve a list of all registered users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200   {array}   models.User
// @Failure      500   {object}  map[string]string
// @Router       /users [get]
func (c *UsersController) GetAll(w http.ResponseWriter, r *http.Request) {

	users, err := c.services.Users.GetAll(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, users)
}

// GetByID retrieves a user by ID
// @Summary      Get user by ID
// @Description  Retrieve details of a specific user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id  path      int  true  "User ID"
// @Success      200      {object}  models.User
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/{user_id} [get]
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

// GetUserTeam retrieves the team associated with a user
// @Summary      Get user's team
// @Description  Retrieve the team information for a specific user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id  path      int  true  "User ID"
// @Success      200      {object}  models.TeamMember
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/{user_id}/teams [get]
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

// Update modifies an existing user
// @Summary      Update user
// @Description  Update details of an existing user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id  path      int          true  "User ID"
// @Param        user     body      models.User  true  "User Update Data"
// @Success      200      {object}  map[string]uint64
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/{user_id} [patch]
func (c *UsersController) Update(w http.ResponseWriter, r *http.Request) {

	userIDToken, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("UserKey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(userIDToken, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
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

// Delete removes a user
// @Summary      Delete user
// @Description  Remove a user from the system
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user_id  path      int  true  "User ID"
// @Success      200      {object}  map[string]uint64
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users/{user_id} [delete]
func (c *UsersController) Delete(w http.ResponseWriter, r *http.Request) {

	userIDToken, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("UserKey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(userIDToken, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
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
		"affected_rows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}
