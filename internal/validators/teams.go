package validators

import (
	"HareID/internal/repository"
	"context"
)

type TeamValidations struct {
	repo repository.Repository
}

func (r *TeamValidations) IsTeamOwner(ctx context.Context, userID, TeamID uint64) (bool, error) {

	team, err := r.repo.Teams.GetByID(ctx, TeamID)
	if err != nil {
		return false, err
	}

	return userID == team.OwnerID, nil
}
