package services

import (
	"HareCRM/internal/models"
	"HareCRM/internal/repository"
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
		Update(ctx context.Context, userID, requestUserID uint64, user models.User) (uint64, error)
		Delete(ctx context.Context, userID, requestUserID uint64) (uint64, error)
	}
}

func NewServices(repository repository.Repository, db *pgxpool.Pool) Services {
	return Services{
		Users: &UserServices{repository, db},
	}
}
