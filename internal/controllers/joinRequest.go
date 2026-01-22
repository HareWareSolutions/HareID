package controllers

import (
	"HareCRM/internal/enums"
	"HareCRM/internal/middleware"
	"HareCRM/internal/responses"
	"HareCRM/internal/services"
	"errors"
	"net/http"
	"strconv"
)

type JoinRequestsController struct {
	services services.Services
}

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
		enums.MEMBER,
		teamID,
		joinRequestData.SenderID,
	)

	data := map[string]interface{}{
		"affectedRows": affectedRows,
		"teamMember":   createdTeamMember,
	}

	responses.JSON(w, http.StatusOK, data)
}

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
