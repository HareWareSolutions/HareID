package models

import (
	"time"

	"HareCRM/internal/enums"
)

type TeamMember struct {
	ID        uint64         `json:"id,omitempty"`
	TeamID    uint64         `json:"team_id,omitempty"`
	UserID    uint64         `json:"user_id,omitempty"`
	Role      enums.TeamRole `json:"role"`
	Name      string         `json:"name,omitempty"`
	TeamName  string         `json:"team_name,omitempty"`
	Email     string         `json:"email,omitempty"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
}
