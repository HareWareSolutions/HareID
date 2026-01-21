package services

import (
	"HareCRM/internal/enums"
	"HareCRM/internal/models"
	"HareCRM/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamMembers struct {
	repo repository.Repository
	db   *pgxpool.Pool
}

func (s *TeamMembers) CreateTeamMember(ctx context.Context, role enums.TeamRole, teamID, userID uint64) (models.TeamMember, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return models.TeamMember{}, err
	}
	defer tx.Rollback(ctx)

	teamMember := models.TeamMember{
		Role:   role,
		TeamID: teamID,
		UserID: userID,
	}

	teamMember, err = s.repo.TeamMembers.Create(ctx, tx, teamMember)
	if err != nil {
		return models.TeamMember{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return models.TeamMember{}, err
	}

	return teamMember, nil
}

func (s *TeamMembers) GetAll(ctx context.Context, teamID uint64) ([]models.TeamMember, error) {

	teamMembers, err := s.repo.TeamMembers.GetAll(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return teamMembers, nil
}

func (s *TeamMembers) GetByUserID(ctx context.Context, userID uint64) (models.TeamMember, error) {

	teamMember, err := s.repo.TeamMembers.GetByUserID(ctx, userID)
	if err != nil {
		return models.TeamMember{}, err
	}

	return teamMember, nil

}
