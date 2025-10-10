package api

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"strconv"
	"fmt"
)

// ArticleHandler struct to handle Article-related requests for future use
type ArticleHandler struct {
}

// NewArticleHandler creates a new instance of ArticleHandler.
func NewArticleHandler() *ArticleHandler {
	return &ArticleHandler{}
}

// HandleGetArticleByID handles the GET request to retrieve a Article by its ID.
func (wh *ArticleHandler) HandleGetArticleByID(w http.ResponseWriter, r *http.Request) {
	
	// retrieves the Article ID from the URL parameters
	paramsArticleID := chi.URLParam(r, "id")
	if paramsArticleID == "" {
		http.NotFound(w, r)
		return
	}

	// checks the validity of the Article ID
	ArticleID, err := strconv.ParseInt(paramsArticleID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Article ID: %d", ArticleID)

}

// HandleCreateArticle handles the POST request to create a new Article.
func (wh *ArticleHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Create Article")
}
