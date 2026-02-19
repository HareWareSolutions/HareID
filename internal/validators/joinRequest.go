package validators

import (
	"HareID/internal/repository"
	"context"
)

type JoinRequestValidations struct {
	repo repository.Repository
}

func (v *JoinRequestValidations) CanSee(ctx context.Context, requestUserID, requestID, teamID uint64) (bool, error) {

	request, err := v.repo.JoinRequests.GetByID(ctx, requestID, teamID)
	if err != nil {
		return false, err
	}

	if request.TeamOwnerID == requestUserID {
		return true, nil
	}

	if request.SenderID == requestUserID {
		return true, nil
	}

	return false, nil
}
