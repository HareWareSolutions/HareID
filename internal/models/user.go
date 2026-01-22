package models

import (
	"HareID/internal/enums"
	"errors"
	"strings"
	"time"
)

type User struct {
	ID        uint64 `json:"id,omitempty"`
	GoogleSub string `json:"google_sub, omitempty"`
	Name      string `json:"name,omitempty"`
	CpfCnpj   string `json:"cpf_cnpj, omnitempty"`
	// Provedor Autenticação - Google - Senha
	StripeCustomerID string             `json:"stripe_customer_id, omitempty"`
	AuthProvider     enums.AuthProvider `json:"auth_provider, omitempty"`
	ConsentTerms     bool               `json:consent_termns,omitempty`
	DataConsent      time.Time          `json:"data_consent,omitempty"`
	CreateDate       time.Time          `json:"create_date,omitempty"`
	UpdateDate       time.Time          `json:"update_date,omitempty"`
}

// Valida os dados do usuário
func (user *User) ValidateUser(step string) error {
	if err := user.ValidateData(step); err != nil {
		return err
	}

	return nil
}

// Valida os dados dos campos do usuário
func (user *User) ValidateData(step string) error {
	if step != "update" {
		if user.GoogleSub == "" {
			return errors.New("google_sub is required")
		}

		if user.Name == "" {
			return errors.New("name is required")
		}

		if user.AuthProvider != 0 && user.AuthProvider != 1 {
			return errors.New("auth_provider is required")
		}
	}

	if step == "login" {
		if !user.ConsentTerms {
			return errors.New("consent_terms is required and cannot be refused")
		}
	}

	return nil
}

// Formata os campos de nome e CPF/CNPJ do usuário
func (user *User) Format() {
	user.Name = strings.TrimSpace(user.Name)
	user.CpfCnpj = strings.TrimSpace(user.CpfCnpj)
}
