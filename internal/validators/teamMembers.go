package validators

import (
	"HareCRM/internal/repository"
	"context"
)

type TeamMemberValidations struct {
	repo repository.Repository
}

func (v *TeamMemberValidations) IsTeamMember(ctx context.Context, userID, teamID uint64) (bool, error) {
	return true, nil
}
