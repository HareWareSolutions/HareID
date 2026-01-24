package main

import (
	"HareID/internal/controllers"
	"HareID/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

func createRouter(controllers controllers.Controller) http.Handler {

	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Origem do seu Angular
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		// Debug: true, // Ative para ver logs de CORS no terminal se der erro
	})

	//Rotas de usuários
	router.Post("/webhook", controllers.Webhook.HandleWebhook)
	router.Post("/login", controllers.Login.Login)
	router.Post("/users", controllers.Users.Create)
	router.Get("/users", middleware.Authenticate(controllers.Users.GetAll))
	router.Get("/users/{user_id}", middleware.Authenticate(controllers.Users.GetByID))
	router.Patch("/users/{user_id}", middleware.Authenticate(controllers.Users.Update))
	router.Delete("/users/{user_id}", middleware.Authenticate(controllers.Users.Delete))

	router.Get("/users/{user_id}/teams", controllers.Users.GetUserTeam)

	//Rotas de Subscriptions
	router.Post("/checkout-session", middleware.Authenticate(controllers.Checkout.CreateSession))
	router.Post("/subscriptions", middleware.Authenticate(controllers.Subscriptions.Create))
	router.Get("/subscriptions", middleware.Authenticate(controllers.Subscriptions.GetAll))
	router.Get("/subscriptions/{subscription_id}", middleware.Authenticate(controllers.Subscriptions.GetBySubscriptionID))
	router.Patch("/subscriptions/{subscription_id}", middleware.Authenticate(controllers.Subscriptions.Update))
	router.Delete("/subscriptions/{subscription_id}", middleware.Authenticate(controllers.Subscriptions.Delete))

	//Rotas de teams
	router.Post("/teams", middleware.Authenticate(controllers.Teams.Create))
	router.Get("/teams", middleware.Authenticate(controllers.Teams.GetAll))
	router.Get("/teams/{team_id}", middleware.Authenticate(controllers.Teams.GetByID))
	router.Patch("/teams/{team_id}", middleware.Authenticate(controllers.Teams.Update))
	router.Delete("/teams/{team_id}", middleware.Authenticate(controllers.Teams.Delete))

	// Rotas de Team Member
	router.Get("/teams/{team_id}/members", middleware.Authenticate(controllers.Teams.GetTeamMembers))

	//Rotas de Join Request
	router.Post("/teams/{team_id}/join", middleware.Authenticate(controllers.JoinRequests.Create))
	router.Get("/teams/{team_id}/join-requests", middleware.Authenticate(controllers.JoinRequests.GetAll))
	router.Get("/teams/{team_id}/join-requests/{request_id}", middleware.Authenticate(controllers.JoinRequests.GetByID))
	router.Delete("/teams/{team_id}/join-requests/{request_id}", middleware.Authenticate(controllers.JoinRequests.Delete))
	router.Patch("/teams/{team_id}/join-requests/{request_id}/accept", middleware.Authenticate(controllers.JoinRequests.Accept))
	router.Patch("/teams/{team_id}/join-requests/{request_id}/reject", middleware.Authenticate(controllers.JoinRequests.Reject))

	//Rotas de Notificacões
	router.Get("/users/{user_id}/notifications", middleware.Authenticate(controllers.Notifications.GetAll))
	router.Get("/users/{user_id}/notifications/{notification_id}", middleware.Authenticate(controllers.Notifications.GetByID))
	router.Delete("/users/{user_id}/notifications/{notification_id}", middleware.Authenticate(controllers.Notifications.Delete))

	handler := cors.Handler(router)

	return handler
}
