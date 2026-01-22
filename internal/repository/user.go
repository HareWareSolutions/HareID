package repository

import (
	"HareID/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func (r UserRepository) Create(ctx context.Context, tx pgx.Tx, user models.User) (models.User, error) {
	query := `
		INSERT INTO users (google_sub, name, cpf_cnpj, stripe_customer_id, auth_provider, consent_terms, data_consent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, create_date
	`

	err := tx.QueryRow(ctx, query,
		user.GoogleSub,
		user.Name,
		user.CpfCnpj,
		user.StripeCustomerID,
		user.AuthProvider,
		user.ConsentTerms,
		user.DataConsent,
	).Scan(&user.ID, &user.CreateDate)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT id, name, cpf_cnpj, stripe_customer_id FROM users
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Name, &user.CpfCnpj, &user.StripeCustomerID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r UserRepository) GetByGoogleSubscription(ctx context.Context, googleSubscription string) (models.User, error) {

	query := `
		SELECT id, google_sub,name, cpf_cnpj, stripe_customer_id, auth_provider, consent_terms, data_consent, create_date
		FROM users
		WHERE google_sub = $1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, googleSubscription).Scan(
		&user.ID,
		&user.Name,
		&user.CpfCnpj,
		&user.StripeCustomerID,
		&user.AuthProvider,
		&user.ConsentTerms,
		&user.DataConsent,
		&user.CreateDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}

func (r UserRepository) GetByID(ctx context.Context, userID uint64) (models.User, error) {
	query := `
		SELECT id, name, cpf_cnpj, stripe_customer_id, auth_provider, consent_terms, data_consent, create_date
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Name,
		&user.CpfCnpj,
		&user.StripeCustomerID,
		&user.AuthProvider,
		&user.ConsentTerms,
		&user.DataConsent,
		&user.CreateDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	return user, nil
}

func (r UserRepository) Update(ctx context.Context, tx pgx.Tx, userID uint64, user models.User) (uint64, error) {

	query := `
		UPDATE users
		SET name = $1, cpf_cnpj = $2, update_date = NOW()
		WHERE id = $3
	`

	result, err := tx.Exec(ctx, query, user.Name, user.CpfCnpj, userID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, errors.New("no user updated")
	}

	return uint64(result.RowsAffected()), nil
}

func (r UserRepository) Delete(ctx context.Context, tx pgx.Tx, userID uint64) (uint64, error) {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := tx.Exec(ctx, query, userID)
	if err != nil {
		return 0, err
	}

	if result.RowsAffected() == 0 {
		return 0, err
	}

	return uint64(result.RowsAffected()), nil
}
