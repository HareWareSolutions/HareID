package repository

import (
	"HareCRM/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Users interface {
		Create(ctx context.Context, tx pgx.Tx, user models.User) (models.User, error)
		GetAll(ctx context.Context) ([]models.User, error)
		GetByGoogleSubscription(ctx context.Context, googleSubscription string) (models.User, error)
		GetByID(ctx context.Context, userID uint64) (models.User, error)
		Update(ctx context.Context, tx pgx.Tx, userID uint64, user models.User) (uint64, error)
		Delete(ctx context.Context, tx pgx.Tx, userID uint64) (uint64, error)
	}
	Teams interface {
		Create(ctx context.Context, tx pgx.Tx, team models.Team) (models.Team, error)
		GetAll(ctx context.Context, tx pgx.Tx) ([]models.Team, error)
		GetByID(ctx context.Context, tx pgx.Tx, teamID uint64) (models.Team, error)
		SearchByOwnerID(ctx context.Context, tx pgx.Tx, userID uint64) (models.Team, error)
		Update(ctx context.Context, tx pgx.Tx, teamID uint64, team models.Team) (uint64, error)
		Delete(ctx context.Context, tx pgx.Tx, teamID uint64) (uint64, error)
	}
	TeamMembers interface {
		Create(ctx context.Context, tx pgx.Tx, teamMember models.TeamMember) (models.TeamMember, error)
		GetTeamMembers(ctx context.Context, tx pgx.Tx, teamID uint64) ([]models.TeamMember, error)
		GetTeamMemberByUserID(ctx context.Context, tx pgx.Tx, userID uint64) (models.TeamMember, error)
	}
	JoinRequests interface {
		Create(ctx context.Context, tx pgx.Tx, joinRequest models.JoinRequest) (models.JoinRequest, error)
		GetAll(ctx context.Context, tx pgx.Tx, teamID uint64) ([]models.JoinRequest, error)
		GetByID(ctx context.Context, tx pgx.Tx, joinRequestID, teamID uint64) (models.JoinRequest, error)
		Delete(ctx context.Context, tx pgx.Tx, requestID, teamID uint64) (uint64, error)
		Accept(ctx context.Context, tx pgx.Tx, userID, teamID, joinRequestID uint64) (uint64, error)
		Reject(ctx context.Context, tx pgx.Tx, userID, teamID, joinRequestID uint64) (uint64, error)
	}
	Notifications interface {
		CreateByJoinRequest(ctx context.Context, tx pgx.Tx, joinRequest models.JoinRequest) (models.Notification, error)
		GetAll(ctx context.Context, tx pgx.Tx, userID uint64) ([]models.Notification, error)
		GetByID(ctx context.Context, tx pgx.Tx, userID, notificationID uint64) (models.Notification, error)
		Delete(ctx context.Context, tx pgx.Tx, userID, notificationID uint64) (uint64, error)
	}
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		Users:         &UserRepository{db: db},
		Teams:         &TeamsRepository{db: db},
		TeamMembers:   &TeamMembersRepository{db: db},
		JoinRequests:  &JoinRequestRepository{db: db},
		Notifications: &NotificationRepository{db: db},
	}
}
