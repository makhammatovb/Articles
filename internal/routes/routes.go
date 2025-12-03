package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/makhammatovb/Articles/internal/app"
)

// SetupRoutes sets up the routes for the application using chi router
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Use(app.Middleware.Authenticate)

	r.Group(func(r chi.Router){
		r.Use(app.Middleware.RequireAuthenticatedUser)

		r.Post("/articles/", app.ArticleHandler.HandleCreateArticle)        // checked
		r.Put("/articles/{id}/", app.ArticleHandler.HandleUpdateArticle)    // checked
		r.Delete("/articles/{id}/", app.ArticleHandler.HandleDeleteArticle) // checked

		r.Put("/users/{id}/", app.UserHandler.HandleUpdateUser)    // checked
		r.Delete("/users/{id}/", app.UserHandler.HandleDeleteUser) // checked
		r.Post("/users/{id}/password-change/", app.UserHandler.HandleUpdatePassword) // checked
		r.Get("/users/{id}", app.UserHandler.HandleGetUserByID)    // checked

		//reviews
		r.Post("/reviews/", app.ReviewHandler.HandleCreateReview)        // checked
		r.Put("/reviews/{id}/", app.ReviewHandler.HandleUpdateReview)    // checked
		r.Delete("/reviews/{id}/", app.ReviewHandler.HandleDeleteReview) // checked

		
	})
	// articles
	r.Get("/health", app.HealthCheck)                                   // checked
	// users
	r.Post("/users/", app.UserHandler.HandleRegisterUser)      			// checked

	r.Get("/articles/{id}", app.ArticleHandler.HandleGetArticleByID)    // checked

	r.Get("/reviews/{id}", app.ReviewHandler.HandleGetReviewByID)    // checked

	// users password update
	r.Post("/users/reset-password-request/", app.TokenHandler.GenerateResetPasswordToken)
	r.Post("/users/reset-password/{token}/", app.TokenHandler.HandleResetPassword)

	// tokens
	r.Post("/tokens/", app.TokenHandler.HandleCreateToken)
	return r
}
