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

func (c *TeamsController) GetAll(w http.ResponseWriter, r *http.Request) {

	teams, err := c.services.Teams.GetAll(r.Context())
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusFound, teams)
}

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
