package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/makhammatovb/Articles/internal/api"
	"github.com/makhammatovb/Articles/internal/store"
	"github.com/makhammatovb/Articles/migrations"
)

// Application struct includes logger and handler from api package
// central container for keeping application-wide dependencies
type Application struct {
	Logger         *log.Logger
	ArticleHandler *api.ArticleHandler
	UserHandler    *api.UserHandler
	ReviewHandler  *api.ReviewHandler
	DB             *sql.DB
}

// NewApplication creates a new instance of Application
// and returns a pointer to it with error
func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	articleStore := store.NewPostgresArticleStore(pgDB)
	userStore := store.NewPostgresUserStore(pgDB)
	reviewStore := store.NewPostgresReviewStore(pgDB)
	// Initialize handlers from api package, creates a new instance of ArticleHandler and returns pointer to it
	ArticleHandler := api.NewArticleHandler(articleStore, logger)
	UserHandler := api.NewUserHandler(userStore, logger)
	ReviewHandler := api.NewReviewHandler(reviewStore, articleStore, logger)
	app := &Application{
		Logger:         logger,
		ArticleHandler: ArticleHandler,
		UserHandler:    UserHandler,
		ReviewHandler:  ReviewHandler,
		DB:             pgDB,
	}
	return app, nil
}

// HealthCheck is a simple handler to check the health of the application
func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Status is available")
}
