package controllers

import (
	"HareID/internal/enums"
	"HareID/internal/middleware"
	"HareID/internal/responses"
	"HareID/internal/services"
	"errors"
	"net/http"
	"strconv"
)

type JoinRequestsController struct {
	services services.Services
}

// Create sends a request to join a team
// @Summary      Request to join team
// @Description  Create a join request for a specific team
// @Tags         join-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id  path      int  true  "Team ID"
// @Success      201      {object}  map[string]any
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /teams/{team_id}/join [post]
func (c *JoinRequestsController) Create(w http.ResponseWriter, r *http.Request) {
	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("User id not found in context"))
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

	newJoinRequest, newNotification, err := c.services.JoinRequests.Create(r.Context(), requestUserID, teamID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]any{
		"newJoinRequest":  newJoinRequest,
		"newNotification": newNotification,
	}

	responses.JSON(w, http.StatusCreated, data)

}

// GetAll retrieves all join requests for a team
// @Summary      Get team join requests
// @Description  Retrieve a list of all join requests for a specific team
// @Tags         join-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id  path      int  true  "Team ID"
// @Success      200      {array}   models.JoinRequest
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /teams/{team_id}/join-requests [get]
func (j *JoinRequestsController) GetAll(w http.ResponseWriter, r *http.Request) {

	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("Userkey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	requests, err := j.services.JoinRequests.GetAll(r.Context(), requestUserID, teamID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, requests)
}

// GetByID retrieves a specific join request
// @Summary      Get join request
// @Description  Retrieve details of a specific join request
// @Tags         join-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id     path      int  true  "Team ID"
// @Param        request_id  path      int  true  "Request ID"
// @Success      200         {object}  models.JoinRequest
// @Failure      400         {object}  map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /teams/{team_id}/join-requests/{request_id} [get]
func (j *JoinRequestsController) GetByID(w http.ResponseWriter, r *http.Request) {

	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("Userkey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	requestID, err := strconv.ParseUint(r.PathValue("request_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	request, err := j.services.JoinRequests.GetByID(r.Context(), requestUserID, teamID, requestID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusOK, request)

}

// Delete cancels a join request
// @Summary      Delete join request
// @Description  Cancel or delete a join request
// @Tags         join-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id     path      int  true  "Team ID"
// @Param        request_id  path      int  true  "Request ID"
// @Success      200         {object}  map[string]uint64
// @Failure      400         {object}  map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /teams/{team_id}/join-requests/{request_id} [delete]
func (j *JoinRequestsController) Delete(w http.ResponseWriter, r *http.Request) {
	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("Userkey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	requestID, err := strconv.ParseUint(r.PathValue("request_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	affectedRows, err := j.services.JoinRequests.Delete(r.Context(), requestUserID, teamID, requestID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]uint64{
		"affectedRows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}

// Accept approves a join request
// @Summary      Accept join request
// @Description  Approve a user's request to join a team
// @Tags         join-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id     path      int  true  "Team ID"
// @Param        request_id  path      int  true  "Request ID"
// @Success      200         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /teams/{team_id}/join-requests/{request_id}/accept [patch]
func (j *JoinRequestsController) Accept(w http.ResponseWriter, r *http.Request) {

	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("Userkey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	requestID, err := strconv.ParseUint(r.PathValue("request_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	affectedRows, err := j.services.JoinRequests.Accept(r.Context(), requestUserID, teamID, requestID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	joinRequestData, err := j.services.JoinRequests.GetByID(r.Context(), requestUserID, teamID, requestID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	createdTeamMember, err := j.services.TeamMembers.Create(
		r.Context(),
		enums.MARKETING_MEMBER,
		teamID,
		joinRequestData.SenderID,
	)

	data := map[string]interface{}{
		"affectedRows": affectedRows,
		"teamMember":   createdTeamMember,
	}

	responses.JSON(w, http.StatusOK, data)
}

// Reject denies a join request
// @Summary      Reject join request
// @Description  Deny a user's request to join a team
// @Tags         join-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        team_id     path      int  true  "Team ID"
// @Param        request_id  path      int  true  "Request ID"
// @Success      200         {object}  map[string]uint64
// @Failure      400         {object}  map[string]string
// @Failure      401         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /teams/{team_id}/join-requests/{request_id}/reject [patch]
func (j *JoinRequestsController) Reject(w http.ResponseWriter, r *http.Request) {
	requestUserIDString, ok := r.Context().Value(middleware.UserKey).(string)
	if !ok {
		responses.Error(w, http.StatusUnauthorized, errors.New("Userkey not found in the request"))
		return
	}

	requestUserID, err := strconv.ParseUint(requestUserIDString, 10, 64)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	teamID, err := strconv.ParseUint(r.PathValue("team_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	requestID, err := strconv.ParseUint(r.PathValue("request_id"), 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	affectedRows, err := j.services.JoinRequests.Reject(r.Context(), requestUserID, teamID, requestID)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]uint64{
		"affectedRows": affectedRows,
	}

	responses.JSON(w, http.StatusOK, data)
}
