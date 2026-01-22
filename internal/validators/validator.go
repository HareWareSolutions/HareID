package validators

import (
	"HareID/internal/repository"
	"context"
)

type Validations struct {
	Users interface {
		CanModify(requestUserID, userID uint64) bool
	}
	Teams interface {
		IsTeamOwner(ctx context.Context, userID, teamID uint64) (bool, error)
	}
	TeamMember interface {
		IsTeamMember(ctx context.Context, userID, teamID uint64) (bool, error)
	}
	JoinRequest interface {
		CanSee(ctx context.Context, requestUserID, teamID uint64) (bool, error)
	}
}

func NewValidator(r repository.Repository) Validations {
	return Validations{
		Users:       &UserValidations{repo: r},
		Teams:       &TeamValidations{repo: r},
		TeamMember:  &TeamMemberValidations{repo: r},
		JoinRequest: &JoinRequestValidations{repo: r},
	}
}
