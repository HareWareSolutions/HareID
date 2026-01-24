package services

import (
	"HareID/internal/enums"
	"HareID/internal/models"
	"HareID/internal/repository"
	"HareID/internal/validators"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Services struct {
	Login interface {
		Login(ctx context.Context, googleSubscription string) (string, error)
	}
	Users interface {
		Create(ctx context.Context, user models.User) (models.User, error)
		GetAll(ctx context.Context) ([]models.User, error)
		GetByID(ctx context.Context, userID uint64) (models.User, error)
		GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (models.User, error)
		Update(ctx context.Context, userID, requestUserID uint64, user models.User) (uint64, error)
		Delete(ctx context.Context, userID, requestUserID uint64) (uint64, error)
	}
	Subscriptions interface {
		Create(ctx context.Context, subscription models.Subscription) (models.Subscription, error)
		GetAll(ctx context.Context) ([]models.Subscription, error)
		GetBySubscriptionID(ctx context.Context, subscriptionID string) (models.Subscription, error)
		Update(ctx context.Context, subscriptionID string, subscription models.Subscription) (uint64, error)
		Delete(ctx context.Context, subscriptionID string) (uint64, error)
		UpsertSubscription(ctx context.Context, subscription models.Subscription) error
	}
	Teams interface {
		Create(ctx context.Context, requestUserID uint64, team models.Team) (models.Team, models.TeamMember, error)
		GetAll(ctx context.Context) ([]models.Team, error)
		GetByID(ctx context.Context, teamID uint64) (models.Team, error)
		GetByOwnerID(ctx context.Context, userID uint64) (models.Team, error)
		Update(ctx context.Context, teamID, requestUserID uint64, team models.Team) (uint64, error)
		Delete(ctx context.Context, teamID, requestUserID uint64) (uint64, error)
		GetOwnerID(ctx context.Context, teamID uint64) (uint64, error)
		CompareUserIDWithTeamOwnerID(ctx context.Context, userID, teamID uint64) error
	}
	TeamMembers interface {
		Create(ctx context.Context, role enums.TeamRole, teamID, userID uint64) (models.TeamMember, error)
		GetAll(ctx context.Context, teamID uint64) ([]models.TeamMember, error)
		GetByUserID(ctx context.Context, userID uint64) (models.TeamMember, error)
	}
	JoinRequests interface {
		Create(ctx context.Context, requestUserID, teamID uint64) (models.JoinRequest, models.Notification, error)
		GetAll(ctx context.Context, requestUserID, teamID uint64) ([]models.JoinRequest, error)
		GetByID(ctx context.Context, requestUserID, teamID, requestID uint64) (models.JoinRequest, error)
		Delete(ctx context.Context, requestUserID, teamID, requestID uint64) (uint64, error)
		Accept(ctx context.Context, requestUserID, teamID, requestID uint64) (uint64, error)
		Reject(ctx context.Context, requestUserID, teamID, requestID uint64) (uint64, error)
	}
	Notifications interface {
		GetAll(ctx context.Context, requestUserID, userID uint64) ([]models.Notification, error)
		GetByID(ctx context.Context, requestUserID, userID, notificationID uint64) (models.Notification, error)
		Delete(ctx context.Context, requestUserID, userID, notificationID uint64) (uint64, error)
	}
	Checkout interface {
		CreateCheckoutSession(ctx context.Context, userID uint64, priceID, successURL, cancelURL string) (string, error)
	}
}

func NewServices(r repository.Repository, v validators.Validations, db *pgxpool.Pool) Services {
	return Services{
		Login:         &LoginServices{repo: r, db: db},
		Users:         &UserServices{repo: r, db: db},
		Subscriptions: &SubscriptionServices{repo: r, db: db},
		Teams:         &TeamServices{repo: r, db: db},
		TeamMembers:   &TeamMembersServices{repo: r, db: db, val: v},
		JoinRequests:  &JoinRequestServices{repo: r, db: db, val: v},
		Notifications: &NotificationServices{repo: r, db: db, val: v},
		Checkout:      &CheckoutServices{},
	}
}
