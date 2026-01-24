package controllers

import (
	"HareID/internal/middleware"
	"HareID/internal/models"
	"HareID/internal/responses"
	"HareID/internal/services"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type TeamsController struct {
	services services.Services
}

// Create creates a new team
// @Summary      Create a new team
// @Description  Create a new team for the authenticated user and add them as owner
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team  body      models.Team  true  "Team Creation Data"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /teams [post]
func (c *TeamsController) Create(w http.ResponseWriter, r *http.Request) {
	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("User ID not found in context"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	log.Printf("Request user id: %d", requestUserID)

	var team models.Team

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	newTeam, teamMember, err := c.services.Teams.Create(r.Context(), requestUserID, team)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]interface{}{
		"team":       newTeam,
		"teamMember": teamMember,
	}

	responses.JSON(w, http.StatusCreated, data)
}

// GetAll retrieves all teams
// @Summary      Get all teams
// @Description  Retrieve a list of all teams
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      302   {array}   models.Team
// @Failure      500   {object}  map[string]string
// @Router       /teams [get]
func (c *TeamsController) GetAll(w http.ResponseWriter, r *http.Request) {

	teams, err := c.services.Teams.GetAll(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusFound, teams)
}

// GetByID retrieves a team by ID
// @Summary      Get team by ID
// @Description  Retrieve details of a specific team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id  path      int  true  "Team ID"
// @Success      302      {object}  models.Team
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /teams/{team_id} [get]
func (c *TeamsController) GetByID(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	teams, err := c.services.Teams.GetByID(r.Context(), teamID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusFound, teams)
}

func (c *TeamsController) GetByOwnerID(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	teams, err := c.services.Teams.GetByOwnerID(r.Context(), userID)

	responses.JSON(w, http.StatusOK, teams)
}

// GetTeamMembers retrieves members of a team
// @Summary      Get team members
// @Description  Retrieve a list of all members in a specific team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id  path      int  true  "Team ID"
// @Success      200      {array}   models.TeamMember
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /teams/{team_id}/members [get]
func (c *TeamsController) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	teamMembers, err := c.services.TeamMembers.GetAll(r.Context(), teamID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, teamMembers)
}

// Update modifies an existing team
// @Summary      Update team
// @Description  Update details of an existing team
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id  path      int          true  "Team ID"
// @Param        team     body      models.Team  true  "Team Update Data"
// @Success      200      {object}  map[string]uint64
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /teams/{team_id} [patch]
func (c *TeamsController) Update(w http.ResponseWriter, r *http.Request) {

	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("Userkey not found in the context"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	var team models.Team

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	affectedRows, err := c.services.Teams.Update(r.Context(), requestUserID, teamID, team)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]uint64{
		"affected_rows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}

// Delete removes a team
// @Summary      Delete team
// @Description  Remove a team from the system
// @Tags         teams
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id  path      int  true  "Team ID"
// @Success      200      {object}  map[string]uint64
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /teams/{team_id} [delete]
func (c *TeamsController) Delete(w http.ResponseWriter, r *http.Request) {

	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("Userkey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	affectedRows, err := c.services.Teams.Delete(r.Context(), requestUserID, teamID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]uint64{
		"affectedRows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}
