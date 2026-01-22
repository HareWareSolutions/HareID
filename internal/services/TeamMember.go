package services

import (
	"HareID/internal/enums"
	"HareID/internal/models"
	"HareID/internal/repository"
	"HareID/internal/validators"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamMembersServices struct {
	repo repository.Repository
	val  validators.Validations
	db   *pgxpool.Pool
}

func (s *TeamMembersServices) Create(ctx context.Context, role enums.TeamRole, teamID, userID uint64) (models.TeamMember, error) {

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

func (s *TeamMembersServices) GetAll(ctx context.Context, teamID uint64) ([]models.TeamMember, error) {

	teamMembers, err := s.repo.TeamMembers.GetAll(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return teamMembers, nil
}

func (s *TeamMembersServices) GetByUserID(ctx context.Context, userID uint64) (models.TeamMember, error) {

	teamMember, err := s.repo.TeamMembers.GetByUserID(ctx, userID)
	if err != nil {
		return models.TeamMember{}, err
	}

	return teamMember, nil

}
