package repository

import (
	"HareCRM/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamMembersRepository struct {
	db *pgxpool.Pool
}

func (r *TeamMembersRepository) Create(ctx context.Context, tx pgx.Tx, teamMember models.TeamMember) (models.TeamMember, error) {

	query := `
		INSERT INTO teammembers (role, team_id, user_id)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	if err := tx.QueryRow(
		ctx,
		query,
		teamMember.Role,
		teamMember.TeamID,
		teamMember.UserID,
	).Scan(
		&teamMember.ID,
		&teamMember.CreatedAt,
	); err != nil {
		return models.TeamMember{}, err
	}

	return teamMember, nil

}

func (r *TeamMembersRepository) GetAll(ctx context.Context, teamID uint64) ([]models.TeamMember, error) {

	query := `
		SELECT tm.id, tm.role, tm.user_id, tm.created_at, u.name FROM teammembers tm
		INNER JOIN users u on u.id = tm.user_id
		INNER JOIN teams t on t.id = tm.team_id
		WHERE team_id = $1
	`

	rows, err := r.db.Query(
		ctx,
		query,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.TeamMember

	for rows.Next() {
		var member models.TeamMember

		if err := rows.Scan(
			&member.ID,
			&member.Role,
			&member.UserID,
			&member.CreatedAt,
			&member.Name,
		); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func (r *TeamMembersRepository) GetByUserID(ctx context.Context, userID uint64) (models.TeamMember, error) {

	query := `
		SELECT tm.id, tm.role, tm.user_id, tm.created_at, u.name, t.name FROM teammembers tm
		INNER JOIN users u on u.id = tm.user_id
		INNER JOIN teams t on t.id = tm.team_id
		WHERE user_id =  $1
	`

	var teamMember models.TeamMember

	if err := r.db.QueryRow(ctx, query, userID).Scan(
		&teamMember.ID,
		&teamMember.Role,
		&teamMember.UserID,
		&teamMember.CreatedAt,
		&teamMember.Name,
		&teamMember.TeamName,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.TeamMember{}, errors.New("team member not found")
		}
		return models.TeamMember{}, err
	}

	return teamMember, nil

}
