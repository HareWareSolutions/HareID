package validators

import (
	"HareID/internal/models"
	"HareID/internal/repository"
	"context"
	"slices"
)

type TeamMemberValidations struct {
	repo repository.Repository
}

func (v *TeamMemberValidations) IsTeamMember(ctx context.Context, userID, teamID uint64) (bool, error) {
	members, err := v.repo.TeamMembers.GetAll(ctx, teamID)
	if err != nil {
		return false, err
	}

	result := slices.ContainsFunc(members, func(member models.TeamMember) bool {
		return member.UserID == userID
	})

	return result, nil
}
