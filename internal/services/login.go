package services

import (
	"HareID/internal/authentication"
	"HareID/internal/repository"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginServices struct {
	repo repository.Repository
	db   *pgxpool.Pool
}

func (ls *LoginServices) Login(ctx context.Context, googleSubscription string) (string, error) {

	user, err := ls.repo.Users.GetByGoogleSubscription(ctx, googleSubscription)
	if err != nil {
		return "", err
	}

	if err = user.ValidateUser("login"); err != nil {
		return "", err
	}

	token, err := authentication.CreateToken(user.GoogleSub, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
