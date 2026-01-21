package services

import (
	"HareCRM/internal/enums"
	"HareCRM/internal/models"
	"HareCRM/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Teams struct {
	repo repository.Repository
	db   *pgxpool.Pool
}

func (s *Teams) CreateTeam(ctx context.Context, requestUserID uint64, team models.Team) (models.Team, models.TeamMember, error) {

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

func (s *Teams) GetAll(ctx context.Context) ([]models.Team, error) {

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

func (ts *Teams) GetByID(ctx context.Context, teamID uint64) (models.Team, error) {

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

func (ts *Teams) UpdateTeam(ctx context.Context, requestUserID, teamID uint64, team models.Team) (uint64, error) {

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

func (ts *Teams) DeleteTeam(ctx context.Context, requestUserID, teamID uint64) (uint64, error) {

	tx, err := ts.db.Begin(ctx)
	if err != nil {
		return 0, nil
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

func (ts *Teams) GetTeamOwnerID(ctx context.Context, teamID uint64) (uint64, error) {

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

func (ts *Teams) CompareUserIDWithTeamOwnerID(ctx context.Context, userID, teamID uint64) error {

	team, err := ts.repo.Teams.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	if team.OwnerID != userID {
		return errors.New("only the owner can delete the team")
	}

	return nil
}
