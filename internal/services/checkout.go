package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
)

type CheckoutServices struct{}

func (s *CheckoutServices) CreateCheckoutSession(ctx context.Context, userID uint64, priceID, successURL, cancelURL string) (string, error) {
	if priceID == "" {
		return "", errors.New("price_id is required")
	}
	if successURL == "" {
		return "", errors.New("success_url is required")
	}
	if cancelURL == "" {
		return "", errors.New("cancel_url is required")
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL:        stripe.String(successURL),
		CancelURL:         stripe.String(cancelURL),
		ClientReferenceID: stripe.String(fmt.Sprintf("%d", userID)),
	}

	sess, err := session.New(params)
	if err != nil {
		return "", err
	}

	return sess.URL, nil
}
