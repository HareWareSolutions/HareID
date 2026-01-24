package repository

import (
	"HareID/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func (r SubscriptionRepository) Create(ctx context.Context, tx pgx.Tx, subscription models.Subscription) (models.Subscription, error) {

	query := `
		INSERT INTO subscriptions (user_id, subscription_id, price_id, status, current_period_end)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, subscription_id
	`

	err := tx.QueryRow(ctx, query,
		subscription.UserID,
		subscription.SubscriptionID,
		subscription.PriceID,
		subscription.Status,
		subscription.CurrentPeriodEnd,
	).Scan(&subscription.ID, &subscription.SubscriptionID)

	if err != nil {
		return models.Subscription{}, err
	}

	return subscription, nil
}

func (r SubscriptionRepository) GetAll(ctx context.Context) ([]models.Subscription, error) {
	query := `
		SELECT id, user_id, subscription_id, price_id, status, current_period_end FROM subscriptions
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []models.Subscription

	for rows.Next() {
		var subscription models.Subscription

		err := rows.Scan(&subscription.ID, &subscription.UserID, &subscription.SubscriptionID, &subscription.PriceID, &subscription.Status, &subscription.CurrentPeriodEnd)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

func (r SubscriptionRepository) GetBySubscriptionID(ctx context.Context, subscriptionID string) (models.Subscription, error) {
	query := `
		SELECT id, user_id, subscription_id, price_id, status, current_period_end FROM subscriptions
		WHERE subscription_id = $1
	`

	var subscription models.Subscription

	err := r.db.QueryRow(ctx, query, subscriptionID).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.SubscriptionID,
		&subscription.PriceID,
		&subscription.Status,
		&subscription.CurrentPeriodEnd,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Subscription{}, errors.New("subscription not found")
		}
		return models.Subscription{}, err
	}

	return subscription, nil
}

func (r SubscriptionRepository) GetByID(ctx context.Context, id uint64) (models.Subscription, error) {
	query := `
		SELECT * from subscriptions WHERE id = $1
	`

	var subscription models.Subscription

	err := r.db.QueryRow(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.SubscriptionID,
		&subscription.PriceID,
		&subscription.Status,
		&subscription.CurrentPeriodEnd,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Subscription{}, errors.New("subscription not found")
		}
		return models.Subscription{}, err
	}

	return subscription, nil
}

func (r SubscriptionRepository) Update(ctx context.Context, tx pgx.Tx, subscriptionID string, subscription models.Subscription) (uint64, error) {

	query := `
		UPDATE subscriptions
		SET price_id = $1, status = $2, current_period_end = $3
		WHERE subscription_id = $4
	`

	result, err := tx.Exec(ctx, query, subscription.PriceID, subscription.Status, subscription.CurrentPeriodEnd)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no subscription updated")
	}

	return uint64(result.RowsAffected()), nil
}

func (r SubscriptionRepository) Delete(ctx context.Context, tx pgx.Tx, subscriptionID string) (uint64, error) {

	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	result, err := tx.Exec(ctx, query, subscriptionID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, err
	}

	return uint64(result.RowsAffected()), nil
}
