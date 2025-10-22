package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/makhammatovb/Articles/internal/app"
)

// SetupRoutes sets up the routes for the application using chi router
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	// articles
	r.Get("/health", app.HealthCheck) // checked
	r.Get("/articles/{id}", app.ArticleHandler.HandleGetArticleByID) // checked
	r.Post("/articles/", app.ArticleHandler.HandleCreateArticle) // checked
	r.Put("/articles/{id}/", app.ArticleHandler.HandleUpdateArticle) // checked
	r.Delete("/articles/{id}/", app.ArticleHandler.HandleDeleteArticle) // checked
	// users
	r.Get("/users/{id}", app.UserHandler.HandleGetUserByID) // checked
	r.Post("/users/", app.UserHandler.HandleCreateUser) // checked
	r.Put("/users/{id}/", app.UserHandler.HandleUpdateUser) // checked
	r.Delete("/users/{id}/", app.UserHandler.HandleDeleteUser) // checked
	return r
}
