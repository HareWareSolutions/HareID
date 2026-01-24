package controllers

import (
	"HareID/internal/enums/subscription"
	"HareID/internal/models"
	"HareID/internal/services"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v79"
	stripeWebhook "github.com/stripe/stripe-go/v79/webhook"
)

type WebhookController struct {
	services services.Services
}

// HandleWebhook processes Stripe webhook events
// @Summary      Stripe Webhook
// @Description  Receive and process Stripe webhook events (e.g. checkout completed, subscription updated)
// @Tags         webhook
// @Accept       json
// @Produce      json
// @Param        Stripe-Signature header string true "Stripe Signature"
// @Success      200  {string}  string "OK"
// @Failure      400  {string}  string "Bad Request"
// @Failure      503  {string}  string "Service Unavailable"
// @Router       /webhook [post]
func (c *WebhookController) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusServiceUnavailable)
		return
	}

	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	event, err := stripeWebhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), endpointSecret)
	if err != nil {
		http.Error(w, "Invalid Signature", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			http.Error(w, "Error parsing webhook JSON", http.StatusBadRequest)
			return
		}

		if session.ClientReferenceID != "" {
			userID, err := strconv.ParseUint(session.ClientReferenceID, 10, 64)
			if err == nil {
				user, err := c.services.Users.GetByID(r.Context(), userID)
				if err == nil {
					user.StripeCustomerID = session.Customer.ID
					c.services.Users.Update(r.Context(), userID, userID, user)
				}
			}
		}

	case "customer.subscription.updated", "customer.subscription.created", "customer.subscription.deleted":
		var stripeSub stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &stripeSub)
		if err != nil {
			http.Error(w, "Error parsing webhook JSON", http.StatusBadRequest)
			return
		}

		if err := c.processSubscriptionEvent(r.Context(), stripeSub); err != nil {
			http.Error(w, "Error processing event: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (c *WebhookController) processSubscriptionEvent(ctx context.Context, stripeSub stripe.Subscription) error {
	status := mapStripeStatusToEnum(string(stripeSub.Status))

	stripeCustomerID := stripeSub.Customer.ID
	user, err := c.services.Users.GetByStripeCustomerID(ctx, stripeCustomerID)
	if err != nil {
		return nil
	}

	sub := models.Subscription{
		UserID:           user.ID,
		SubscriptionID:   stripeSub.ID,
		PriceID:          stripeSub.Items.Data[0].Price.ID,
		Status:           status,
		CurrentPeriodEnd: time.Unix(stripeSub.CurrentPeriodEnd, 0),
	}

	return c.services.Subscriptions.UpsertSubscription(ctx, sub)
}

func mapStripeStatusToEnum(status string) subscription.Subscription {
	switch strings.ToLower(status) {
	case "active":
		return subscription.ACTIVE
	case "past_due":
		return subscription.PAST_DUE
	case "unpaid":
		return subscription.UNPAID
	case "canceled":
		return subscription.CANCELED
	case "incomplete":
		return subscription.INCOMPLETE
	case "incomplete_expired":
		return subscription.INCOMPLETE_EXPIRED
	case "trialing":
		return subscription.TRIALING
	default:
		return subscription.UNKNOWN
	}
}
