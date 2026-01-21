package validators

import (
	"HareCRM/internal/repository"
	"context"
)

type JoinRequestValidations struct {
	repo repository.Repository
}

func (v *JoinRequestValidations) CanSee(ctx context.Context, userID, teamID uint64) (bool, error) {
	team, err := v.repo.Teams.GetByID(ctx, teamID)
	if err != nil {
		return false, err
	}

	return userID == team.OwnerID, nil
}
