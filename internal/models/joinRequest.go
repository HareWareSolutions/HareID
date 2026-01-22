package models

import (
	"HareID/internal/enums"

	"github.com/lib/pq"
)

type JoinRequest struct {
	ID          uint64       `json:"id,omitempty"`
	TeamID      uint64       `json:"team_id,omitempty"`
	TeamOwnerID uint64       `json:"team_owner_id,omitempty"`
	SenderID    uint64       `json:"sender_id,omitempty"`
	Status      enums.Status `json:"status"`
	DecisionAt  pq.NullTime  `json:"decision_at"`
	DecisionBy  *uint64      `json:"decision_by,omitempty"`
}
