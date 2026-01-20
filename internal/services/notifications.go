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

func NewNotificationServices(db *pgxpool.Pool) *Notifications {
	return &Notifications{db: db}
}

func (s *Notifications) GetAll(ctx context.Context, userID uint) ([]models.Notification, error) {
	notifications, err := s.repo.Notifications.GetAll(ctx, s.db, uint64(userID))
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *Notifications) GetByID(ctx context.Context, userID, notificationID uint64) (models.Notification, error) {

	notification, err := s.repo.Notifications.GetByID(ctx, s.db, userID, notificationID)
	if err != nil {
		return models.Notification{}, err
	}

	return notification, nil
}

func (s *Notifications) Delete(ctx context.Context, userID, notificationID uint64) (uint64, error) {

	affectedRows, err := s.repo.Notifications.Delete(ctx, s.db, userID, notificationID)
	if err != nil {
		return 0, err
	}

	return affectedRows, nil
}
