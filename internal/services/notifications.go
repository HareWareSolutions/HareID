package services

import (
	"HareCRM/internal/models"
	"HareCRM/internal/repository"
	"HareCRM/internal/validators"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationServices struct {
	repo repository.Repository
	val  validators.Validations
	db   *pgxpool.Pool
}

func (s *NotificationServices) GetAll(ctx context.Context, requestUserID, userID uint64) ([]models.Notification, error) {

	if requestUserID != userID {
		return nil, errors.New("you can only see your own notifications")
	}

	notifications, err := s.repo.Notifications.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *NotificationServices) GetByID(ctx context.Context, requestUserID, userID, notificationID uint64) (models.Notification, error) {

	if requestUserID != userID {
		return models.Notification{}, errors.New("you can only see your own notifications")
	}

	notification, err := s.repo.Notifications.GetByID(ctx, userID, notificationID)
	if err != nil {
		return models.Notification{}, err
	}

	return notification, nil
}

func (s *NotificationServices) Delete(ctx context.Context, requestUserID, userID, notificationID uint64) (uint64, error) {

	if requestUserID != userID {
		return 0, errors.New("you can only delete your own notifications")
	}

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
