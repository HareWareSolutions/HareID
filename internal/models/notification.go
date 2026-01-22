package models

import (
	"HareID/internal/enums"

	"github.com/lib/pq"
)

// ID uint64
// SenderID, receiverID, Type, ReferenceID, Seen, CreatedAt

type Notification struct {
	ID          uint64                 `json:"id,omitempty"`
	SenderID    uint64                 `json:"sender_id,omitempty"`
	ReceiverID  uint64                 `json:"receiver_id,omitempty"`
	Type        enums.NotificationType `json:"notification_type,omitempty"`
	ReferenceID uint64                 `json:"reference_id,omitempty"`
	Seen        bool                   `json:"seen,omitempty"`
	CreatedAt   pq.NullTime            `json:"created_at,omitempty"`
}
