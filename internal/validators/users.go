package validators

import "HareCRM/internal/repository"

type UserValidations struct {
	repo repository.Repository
}

func (v *UserValidations) CanModify(requestUserID, userID uint64) bool {
	return requestUserID == userID
}
