package routes

import (
	"github.com/makhammatovb/Articles/internal/app"
	"github.com/go-chi/chi/v5"
)

// SetupRoutes sets up the routes for the application using chi router
func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/articles/{id}", app.ArticleHandler.HandleGetArticleByID)
	r.Post("/articles/", app.ArticleHandler.HandleCreateArticle)
	r.Put("/articles/{id}/", app.ArticleHandler.HandleUpdateArticle)
	r.Delete("/articles/{id}/", app.ArticleHandler.HandleDeleteArticle)
	return r
}
