package services

import (
	"HareID/internal/models"
	"HareID/internal/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionServices struct {
	repo repository.Repository
	db   *pgxpool.Pool
}

func (s *SubscriptionServices) Create(ctx context.Context, subscription models.Subscription) (models.Subscription, error) {
	if subscription.SubscriptionID == "" {
		return models.Subscription{}, errors.New("subscription_id is required")
	}
	
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return models.Subscription{}, err
	}
	defer tx.Rollback(ctx)

	createdSubscription, err := s.repo.Subscriptions.Create(ctx, tx, subscription)
	if err != nil {
		return models.Subscription{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Subscription{}, err
	}

	return createdSubscription, nil
}

func (s *SubscriptionServices) GetAll(ctx context.Context) ([]models.Subscription, error) {

	subscriptions, err := s.repo.Subscriptions.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (s *SubscriptionServices) GetBySubscriptionID(ctx context.Context, subscriptionID string) (models.Subscription, error) {

	subscription, err := s.repo.Subscriptions.GetBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return models.Subscription{}, err
	}

	return subscription, nil
}

func (s *SubscriptionServices) Update(ctx context.Context, subscriptionID string, subscription models.Subscription) (uint64, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	affectedRows, err := s.repo.Subscriptions.Update(ctx, tx, subscriptionID, subscription)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (s *SubscriptionServices) Delete(ctx context.Context, subscriptionID string) (uint64, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, nil
	}
	defer tx.Rollback(ctx)

	affectedRows, err := s.repo.Subscriptions.Delete(ctx, tx, subscriptionID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return affectedRows, nil
}

func (s *SubscriptionServices) UpsertSubscription(ctx context.Context, subscription models.Subscription) error {
	// Tenta buscar a assinatura existente
	existing, err := s.repo.Subscriptions.GetBySubscriptionID(ctx, subscription.SubscriptionID)
	if err == nil {
		// Se existir, atualiza
		_, err := s.Update(ctx, existing.SubscriptionID, subscription)
		return err
	}

	// Se não existir, cria uma nova
	_, err = s.Create(ctx, subscription)
	return err
}
