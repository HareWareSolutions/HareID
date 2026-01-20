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

func NewTeamMembersServices(db *pgxpool.Pool) *TeamMembers {
	return &TeamMembers{db: db}
}

func (tms *TeamMembers) CreateTeamMember(ctx context.Context, role enums.TeamRole, teamID, userID uint64) (models.TeamMember, error) {

	tx, err := tms.db.Begin(ctx)
	if err != nil {
		return models.TeamMember{}, err
	}
	defer tx.Rollback(ctx)

	teamMember := models.TeamMember{
		Role:   role,
		TeamID: teamID,
		UserID: userID,
	}

	teamMember, err = tms.repo.TeamsMembers.Create(ctx, tx, teamMember)
	if err != nil {
		return models.TeamMember{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return models.TeamMember{}, err
	}

	return teamMember, nil
}

func (tms *TeamMembers) GetTeamMembers(ctx context.Context, teamID uint64) ([]models.TeamMember, error) {

	teamMembers, err := tms.repo.TeamsMembers.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return teamMembers, nil
}

func (tms *TeamMembers) GetTeamMemberByUserID(ctx context.Context, userID uint64) (models.TeamMember, error) {

	teamMember, err := tms.repo.TeamsMembers.GetByID(ctx, userID)
	if err != nil {
		return models.TeamMember{}, err
	}

	return teamMember, nil

}
