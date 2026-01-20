package repository

import (
	"HareCRM/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JoinRequestRepository struct {
	db *pgxpool.Pool
}

func (repository *JoinRequestRepository) Create(ctx context.Context, tx pgx.Tx, joinRequest models.JoinRequest) (models.JoinRequest, error) {

	query := `
		INSERT INTO teamjoinrequests (team_id, team_owner_id, sender_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	if err := tx.QueryRow(
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

func (repository *JoinRequestRepository) GetJoinRequest(ctx context.Context, tx pgx.Tx, teamID uint64) ([]models.JoinRequest, error) {

	query := `
		SELECT 	id, team_id, team_owner_id, sender_id, status, decision_at, decision_by 
		FROM teamjoinrequests
		WHERE team_id = $1
	`

	rows, err := tx.Query(ctx, query, teamID)
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

func (repository *JoinRequestRepository) GetJoinRequestByID(ctx context.Context, tx pgx.Tx, joinRequestID, teamID uint64) (models.JoinRequest, error) {

	query := `
		SELECT 	id, team_id, team_owner_id, sender_id, status, decision_at, decision_by 
		FROM teamjoinrequests
		WHERE id = $1 AND team_id = $2
	`

	var request models.JoinRequest

	if err := tx.QueryRow(
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

func (repository *JoinRequestRepository) Delete(ctx context.Context, tx pgx.Tx, requestID, teamID uint64) (uint64, error) {

	query := `
		DELETE FROM teamjoinrequests
		WHERE id = $1 AND team_id = $2
 	`

	result, err := tx.Exec(ctx, query, requestID, teamID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no team deleted")
	}

	return uint64(result.RowsAffected()), nil

}

func (repository *JoinRequestRepository) Accept(ctx context.Context, tx pgx.Tx, userID, teamID, joinRequestID uint64) (uint64, error) {
	query := `
		UPDATE teamjoinrequests SET status = 1, decision_at = NOW(), decision_by = $1 WHERE id = $2 AND team_id = $3
	`

	result, err := tx.Exec(ctx, query, userID, joinRequestID, teamID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no request accepted")
	}

	return uint64(result.RowsAffected()), nil
}

func (repository *JoinRequestRepository) Reject(ctx context.Context, tx pgx.Tx, userID, teamID, joinRequestID uint64) (uint64, error) {
	query := `
		UPDATE teamjoinrequests SET status = 2, decision_at = NOW(), decision_by = $1 WHERE id = $2 AND team_id = $3
	`

	result, err := tx.Exec(ctx, query, userID, joinRequestID, teamID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no request accepted")
	}

	return uint64(result.RowsAffected()), nil
}

func (repository *JoinRequestRepository) Update(ctx context.Context, tx pgx.Tx, joinRequestID uint64, joinRequest models.JoinRequest) (uint64, error) {
	query := `
		UPDATE teamjoinrequests
		SET status = $1
		WHERE id = $2
	`

	result, err := tx.Exec(ctx, query, joinRequest.Status, joinRequestID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no join request updated")
	}

	return uint64(result.RowsAffected()), nil
}
