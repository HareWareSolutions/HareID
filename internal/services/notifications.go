package services

import (
	"HareCRM/internal/models"
	"HareCRM/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Notifications struct {
	repo repository.Repository
	db   *pgxpool.Pool
}

func (s *Notifications) GetAll(ctx context.Context, userID uint) ([]models.Notification, error) {
	notifications, err := s.repo.Notifications.GetAll(ctx, uint64(userID))
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *Notifications) GetByID(ctx context.Context, userID, notificationID uint64) (models.Notification, error) {

	notification, err := s.repo.Notifications.GetByID(ctx, userID, notificationID)
	if err != nil {
		return models.Notification{}, err
	}

	return notification, nil
}

func (s *Notifications) Delete(ctx context.Context, userID, notificationID uint64) (uint64, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	affectedRows, err := s.repo.Notifications.Delete(ctx, tx, userID, notificationID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}
