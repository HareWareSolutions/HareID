package services

import (
	"HareID/internal/enums"
	"HareID/internal/models"
	"HareID/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamServices struct {
	repo repository.Repository
	db   *pgxpool.Pool
}

func (s *TeamServices) Create(ctx context.Context, requestUserID uint64, team models.Team) (models.Team, models.TeamMember, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return models.Team{}, models.TeamMember{}, err
	}
	defer tx.Rollback(ctx)

	team.OwnerID = requestUserID

	if err := team.ValidateTeam("creation"); err != nil {
		return models.Team{}, models.TeamMember{}, err
	}

	team, err = s.repo.Teams.Create(ctx, tx, team)
	if err != nil {
		return models.Team{}, models.TeamMember{}, err
	}

	teamMember := models.TeamMember{
		TeamID: team.ID,
		UserID: team.OwnerID,
		Role:   enums.OWNER,
	}

	teamMember, err = s.repo.TeamMembers.Create(ctx, tx, teamMember)
	if err != nil {
		return models.Team{}, models.TeamMember{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Team{}, models.TeamMember{}, err
	}

	return team, models.TeamMember{}, nil
}

func (s *TeamServices) GetAll(ctx context.Context) ([]models.Team, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	teams, err := s.repo.Teams.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(teams) < 1 {
		return nil, errors.New("no team found")
	}

	return teams, nil
}

func (ts *TeamServices) GetByID(ctx context.Context, teamID uint64) (models.Team, error) {

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return models.Team{}, err
	}
	defer tx.Rollback(ctx)

	team, err := ts.repo.Teams.GetByID(ctx, teamID)
	if err != nil {
		return models.Team{}, err
	}

	return team, nil
}

func (ts *TeamServices) GetByOwnerID(ctx context.Context, userID uint64) (models.Team, error) {
	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return models.Team{}, err
	}
	defer tx.Rollback(ctx)

	team, err := ts.repo.Teams.SearchByOwnerID(ctx, userID)
	if err != nil {
		return models.Team{}, err
	}

	return team, nil
}

func (ts *TeamServices) Update(ctx context.Context, teamID, requestUserID uint64, team models.Team) (uint64, error) {

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	if err := team.ValidateTeam("update"); err != nil {
		return 0, err
	}

	affectedRows, err := ts.repo.Teams.Update(ctx, tx, teamID, team)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (ts *TeamServices) Delete(ctx context.Context, teamID, requestUserID uint64) (uint64, error) {

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	if err := ts.CompareUserIDWithTeamOwnerID(ctx, requestUserID, teamID); err != nil {
		return 0, err
	}

	affectedRows, err := ts.repo.Teams.Delete(ctx, tx, teamID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (ts *TeamServices) GetOwnerID(ctx context.Context, teamID uint64) (uint64, error) {

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	team, err := ts.repo.Teams.GetByID(ctx, teamID)
	if err != nil {
		return 0, err
	}

	return team.OwnerID, nil
}

func (ts *TeamServices) CompareUserIDWithTeamOwnerID(ctx context.Context, userID, teamID uint64) error {

	team, err := ts.repo.Teams.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	if team.OwnerID != userID {
		return errors.New("only the owner can delete the team")
	}

	return nil
}
