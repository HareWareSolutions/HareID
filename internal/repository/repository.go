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
}

func NewRepository(db *pgxpool.Pool) Repository {
	return Repository{
		Users: &UserRepository{db: db},
	}
}
