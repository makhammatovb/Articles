package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/makhammatovb/Articles/internal/api"
)

// Application struct includes logger and handler from api package
// central container for keeping application-wide dependencies
type Application struct {
	Logger         *log.Logger
	ArticleHandler *api.ArticleHandler
}

// NewApplication creates a new instance of Application
// and returns a pointer to it with error
func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize handlers from api package, creates a new instance of ArticleHandler and returns pointer to it
	ArticleHandler := api.NewArticleHandler()
	app := &Application{
		Logger:         logger,
		ArticleHandler: ArticleHandler,
	}
	return app, nil
}

// HealthCheck is a simple handler to check the health of the application
func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available")
}
