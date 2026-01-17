package models

import (
	"errors"
	"time"
)

type Team struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Domain    string    `json:"domain,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (team *Team) ValidadeTeam(step string) error {
	if err := team.ValidateData(step); err != nil {
		return err
	}
	return nil
}

func (team *Team) ValidateData(step string) error {
	if team.Name == "" {
		return errors.New("name is required")
	}
	if team.Domain == "" {
		return errors.New("domain is required")
	}
	return nil
}
