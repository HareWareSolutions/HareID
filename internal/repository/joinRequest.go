package repository

import (
	"HareCRM/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type TeamJoinRequest struct {
}

func NewTeamJoinRequestRepository() *TeamJoinRequest {
	return &TeamJoinRequest{}
}

func (repository *TeamJoinRequest) Create(ctx context.Context, db DBTX, joinRequest models.JoinRequest) (models.JoinRequest, error) {

	query := `
		INSERT INTO teamjoinrequests (team_id, team_owner_id, sender_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	if err := db.QueryRow(
		ctx,
		query,
		joinRequest.TeamID,
		joinRequest.TeamOwnerID,
		joinRequest.SenderID,
		joinRequest.Status,
	).Scan(
		&joinRequest.ID,
	); err != nil {
		return models.JoinRequest{}, err
	}

	return joinRequest, nil
}

func (repository *TeamJoinRequest) GetJoinRequest(ctx context.Context, db DBTX, teamID uint64) ([]models.JoinRequest, error) {

	query := `
		SELECT 	id, team_id, team_owner_id, sender_id, status, decision_at, decision_by 
		FROM teamjoinrequests
		WHERE team_id = $1
	`

	rows, err := db.Query(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.JoinRequest

	for rows.Next() {
		var request models.JoinRequest

		if err := rows.Scan(
			&request.ID,
			&request.TeamID,
			&request.TeamOwnerID,
			&request.SenderID,
			&request.Status,
			&request.DecisionAt,
			&request.DecisionBy,
		); err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (repository *TeamJoinRequest) GetJoinRequestByID(ctx context.Context, db DBTX, joinRequestID, teamID uint64) (models.JoinRequest, error) {

	query := `
		SELECT 	id, team_id, team_owner_id, sender_id, status, decision_at, decision_by 
		FROM teamjoinrequests
		WHERE id = $1 AND team_id = $2
	`

	var request models.JoinRequest

	if err := db.QueryRow(
		ctx,
		query,
		joinRequestID,
		teamID,
	).Scan(
		&request.ID,
		&request.TeamID,
		&request.TeamOwnerID,
		&request.SenderID,
		&request.Status,
		&request.DecisionAt,
		&request.DecisionBy,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.JoinRequest{}, errors.New("join request not found")
		}
		return models.JoinRequest{}, err
	}

	return request, nil
}

func (repository *TeamJoinRequest) Delete(ctx context.Context, db DBTX, requestID, teamID uint64) (uint64, error) {

	query := `
		DELETE FROM teamjoinrequests
		WHERE id = $1 AND team_id = $2
 	`

	result, err := db.Exec(ctx, query, requestID, teamID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no team deleted")
	}

	return uint64(result.RowsAffected()), nil

}

func (repository *TeamJoinRequest) Accept(ctx context.Context, db DBTX, userID, teamID, joinRequestID uint64) (uint64, error) {
	query := `
		UPDATE teamjoinrequests SET status = 1, decision_at = NOW(), decision_by = $1 WHERE id = $2 AND team_id = $3
	`

	result, err := db.Exec(ctx, query, userID, joinRequestID, teamID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no request accepted")
	}

	return uint64(result.RowsAffected()), nil
}

func (repository *TeamJoinRequest) Reject(ctx context.Context, db DBTX, userID, teamID, joinRequestID uint64) (uint64, error) {
	query := `
		UPDATE teamjoinrequests SET status = 2, decision_at = NOW(), decision_by = $1 WHERE id = $2 AND team_id = $3
	`

	result, err := db.Exec(ctx, query, userID, joinRequestID, teamID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no request accepted")
	}

	return uint64(result.RowsAffected()), nil
}
