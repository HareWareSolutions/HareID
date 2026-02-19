package controllers

import (
	"HareID/internal/services"
	"net/http"
)

type Controller struct {
	Login interface {
		Login(http.ResponseWriter, *http.Request)
	}
	Users interface {
		Create(http.ResponseWriter, *http.Request)
		GetAll(http.ResponseWriter, *http.Request)
		GetByID(http.ResponseWriter, *http.Request)
		GetUserTeam(http.ResponseWriter, *http.Request)
		Update(http.ResponseWriter, *http.Request)
		Delete(http.ResponseWriter, *http.Request)
	}
	Subscriptions interface {
		Create(http.ResponseWriter, *http.Request)
		GetAll(http.ResponseWriter, *http.Request)
		GetBySubscriptionID(http.ResponseWriter, *http.Request)
		Update(http.ResponseWriter, *http.Request)
		Delete(http.ResponseWriter, *http.Request)
	}
	Teams interface {
		Create(http.ResponseWriter, *http.Request)
		GetAll(http.ResponseWriter, *http.Request)
		GetByID(http.ResponseWriter, *http.Request)
		GetByOwnerID(http.ResponseWriter, *http.Request)
		GetTeamMembers(http.ResponseWriter, *http.Request)
		Update(http.ResponseWriter, *http.Request)
		Delete(http.ResponseWriter, *http.Request)
	}
	TeamMembers interface {
		GetAll(http.ResponseWriter, *http.Request)
		GetByID(http.ResponseWriter, *http.Request)
		GetByName(http.ResponseWriter, *http.Request)
		Delete(http.ResponseWriter, *http.Request)
	}
	JoinRequests interface {
		Create(http.ResponseWriter, *http.Request)
		GetAll(http.ResponseWriter, *http.Request)
		GetByID(http.ResponseWriter, *http.Request)
		Delete(http.ResponseWriter, *http.Request)
		Accept(http.ResponseWriter, *http.Request)
		Reject(http.ResponseWriter, *http.Request)
	}
	Notifications interface {
		GetAll(http.ResponseWriter, *http.Request)
		GetByID(http.ResponseWriter, *http.Request)
		Delete(http.ResponseWriter, *http.Request)
	}
	Webhook interface {
		HandleWebhook(http.ResponseWriter, *http.Request)
	}
	Checkout interface {
		CreateSession(http.ResponseWriter, *http.Request)
	}
}

func NewControllers(s services.Services) Controller {
	return Controller{
		Login:         &LoginController{services: s},
		Subscriptions: &SubscriptionsController{services: s},
		Users:         &UsersController{services: s},
		Teams:         &TeamsController{services: s},
		JoinRequests:  &JoinRequestsController{services: s},
		Notifications: &NotificationsController{services: s},
		Webhook:       &WebhookController{services: s},
		Checkout:      &CheckoutController{services: s},
	}
}
