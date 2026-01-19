package repository

import (
	"HareCRM/internal/enums"
	"HareCRM/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Notifications struct {
}

func NewNotificationsRepository() *Notifications {
	return &Notifications{}
}

func (repository *Notifications) CreateByJoinRequest(ctx context.Context, db DBTX, joinRequest models.JoinRequest) (models.Notification, error) {

	notification := models.Notification{
		SenderID:    joinRequest.SenderID,
		ReceiverID:  joinRequest.TeamOwnerID,
		Type:        enums.JOIN_REQUEST, //TIPO: Join Request
		ReferenceID: joinRequest.ID,
		Seen:        false,
	}

	query := `
		INSERT INTO notifications(sender_id, receiver_id, type, reference_id, seen)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	if err := db.QueryRow(
		ctx,
		query,
		notification.SenderID,
		notification.ReceiverID,
		notification.Type,
		notification.ReferenceID,
		notification.Seen,
	).Scan(
		&notification.ID,
		&notification.CreatedAt,
	); err != nil {
		return models.Notification{}, err
	}

	return notification, nil
}

func (repository *Notifications) GetAll(ctx context.Context, db DBTX, userID uint64) ([]models.Notification, error) {

	query := `
		SELECT id, sender_id, receiver_id, type, reference_id, seen, created_at
		FROM notifications
		WHERE receiver_id = $1
	`

	rows, err := db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification

	for rows.Next() {

		var notification models.Notification

		if err = rows.Scan(
			&notification.ID,
			&notification.SenderID,
			&notification.ReceiverID,
			&notification.Type,
			&notification.ReferenceID,
			&notification.Seen,
			&notification.CreatedAt,
		); err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	if len(notifications) < 1 {
		return nil, errors.New("no notifications found")
	}

	return notifications, nil
}

func (repository *Notifications) GetByID(ctx context.Context, db DBTX, userID, notificationID uint64) (models.Notification, error) {

	query := `
		SELECT id, sender_id, receiver_id, type, reference_id, seen, created_at
		FROM notifications
		WHERE id = $1 AND receiver_id = $2
	`

	var notification models.Notification

	if err := db.QueryRow(
		ctx,
		query,
		notificationID,
		userID,
	).Scan(
		&notification.ID,
		&notification.SenderID,
		&notification.ReceiverID,
		&notification.Type,
		&notification.ReferenceID,
		&notification.Seen,
		&notification.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Notification{}, errors.New("notification not found")
		}
		return models.Notification{}, err
	}

	return notification, nil
}

func (repository *Notifications) Delete(ctx context.Context, db DBTX, userID, notificationID uint64) (uint64, error) {
	query := `
		DELETE FROM notifications
		WHERE id = $1 AND receiver_id = $2
	`

	result, err := db.Exec(ctx, query, notificationID, userID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() < 1 {
		return 0, errors.New("no notification deleted")
	}

	return uint64(result.RowsAffected()), nil
}
