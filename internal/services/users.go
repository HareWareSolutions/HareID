package services

import (
	"HareID/internal/models"
	"HareID/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserServices struct {
	repo repository.Repository
	db   *pgxpool.Pool
}

func (s *UserServices) Create(ctx context.Context, user models.User) (models.User, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer tx.Rollback(ctx)

	if err := user.ValidateUser("create"); err != nil {
		return models.User{}, err
	}

	createdUser, err := s.repo.Users.Create(ctx, tx, user)
	if err != nil {
		return models.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func (s *UserServices) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.Users.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserServices) GetByID(ctx context.Context, userID uint64) (models.User, error) {

	user, err := s.repo.Users.GetByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserServices) GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (models.User, error) {

	user, err := s.repo.Users.GetByStripeCustomerID(ctx, stripeCustomerID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserServices) Update(ctx context.Context, userID, requestUserID uint64, user models.User) (uint64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	if userID != requestUserID {
		tx.Rollback(ctx)
		return 0, errors.New("Only the owner can update the user")
	}

	if err := user.ValidateUser("update"); err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	affectedRows, err := s.repo.Users.Update(ctx, tx, userID, user)
	if err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return 0, err
	}

	return affectedRows, nil
}


func (s *UserServices) Delete(ctx context.Context, userID, requestUserID uint64) (uint64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback(ctx)

	if userID != requestUserID {
		return 0, errors.New("Only the owner can delete the user")
	}

	affectedRows, err := s.repo.Users.Delete(ctx, tx, userID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil

}
