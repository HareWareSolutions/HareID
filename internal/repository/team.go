package repository

import (
	"HareID/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamsRepository struct {
	db *pgxpool.Pool
}

func (r *TeamsRepository) Create(ctx context.Context, tx pgx.Tx, team models.Team) (models.Team, error) {

	query := `
		INSERT INTO teams (name, domain, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, domain, created_at
	`

	if err := tx.QueryRow(ctx, query, team.Name, team.Domain, team.OwnerID).Scan(
		&team.ID,
		&team.Name,
		&team.Domain,
		&team.CreatedAt,
	); err != nil {
		return models.Team{}, err
	}

	return team, nil
}

func (r *TeamsRepository) GetAll(ctx context.Context) ([]models.Team, error) {

	query := `
		SELECT id, name, domain, owner_id, created_at, updated_at
		FROM teams
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team

	for rows.Next() {
		var team models.Team

		if err = rows.Scan(
			&team.ID,
			&team.Name,
			&team.Domain,
			&team.OwnerID,
			&team.CreatedAt,
			&team.UpdatedAt,
		); err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamsRepository) GetByID(ctx context.Context, teamID uint64) (models.Team, error) {

	query := `
		SELECT id, name, domain, owner_id, created_at, updated_at
		FROM teams
		WHERE id = $1
	`

	var team models.Team

	err := r.db.QueryRow(ctx, query, teamID).Scan(
		&team.ID,
		&team.Name,
		&team.Domain,
		&team.OwnerID,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Team{}, errors.New("team not found")
		}
		return models.Team{}, err
	}

	return team, nil
}

func (r *TeamsRepository) SearchByOwnerID(ctx context.Context, userID uint64) (models.Team, error) {

	query := `
		SELECT id, name, owner_id, created_at, updated_at
		FROM teams
		WHERE id = $1
	`

	var team models.Team

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&team.ID,
		&team.Name,
		&team.OwnerID,
		&team.CreatedAt,
		&team.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Team{}, errors.New("team not found")
		}
		return models.Team{}, err
	}

	return team, nil
}

func (r *TeamsRepository) Update(ctx context.Context, tx pgx.Tx, teamID uint64, team models.Team) (uint64, error) {

	query := `
		UPDATE teams
		SET name = $1, domain = $2, updated_at = NOW()
		WHERE id = $3
	`

	result, err := tx.Exec(ctx, query, team.Name, team.Domain, teamID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no team updated")
	}

	return uint64(result.RowsAffected()), nil
}

func (r *TeamsRepository) Delete(ctx context.Context, tx pgx.Tx, teamID uint64) (uint64, error) {

	query := `
		DELETE FROM teams
		WHERE id = $1
	`

	result, err := tx.Exec(ctx, query, teamID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no team deleted")
	}

	return uint64(result.RowsAffected()), nil
}
