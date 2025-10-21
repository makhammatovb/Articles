package api

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/makhammatovb/Articles/internal/store"
	"github.com/makhammatovb/Articles/internal/utils"
)

// ArticleHandler struct to handle Article-related requests for future use
type ArticleHandler struct {
	articleStore store.ArticleStore
	logger       *log.Logger
}

// NewArticleHandler creates a new instance of ArticleHandler.
func NewArticleHandler(articleStore store.ArticleStore, logger *log.Logger) *ArticleHandler {
	return &ArticleHandler{
		articleStore: articleStore,
		logger:       logger,
	}
}

// HandleGetArticleByID handles the GET request to retrieve a Article by its ID.
// wh *ArticleHandler is the receiver
// *ArticleHandler means the method operates on a pointer to ArticleHandler
// w http.ResponseWriter is used to write the response back to the client
// r *http.Request represents the incoming HTTP request
func (ah *ArticleHandler) HandleGetArticleByID(w http.ResponseWriter, r *http.Request) {

	// retrieves the article ID from the URL parameters
	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Println("Error reading article ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid article ID"})
		return
	}
	article, err := ah.articleStore.GetArticleByID(articleID)
	if err != nil {
		ah.logger.Println("Error getting article by ID:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"article": article})
}

// HandleCreateArticle handles the POST request to create a new article.
func (ah *ArticleHandler) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	var article store.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		ah.logger.Println("Decoding error:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	createdArticle, err := ah.articleStore.CreateArticle(&article)
	if err != nil {
		ah.logger.Println("Error creating article:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"article": createdArticle})
}

func (ah *ArticleHandler) HandleUpdateArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Println("Error reading article ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid article ID"})
		return
	}
	existingArticle, err := ah.articleStore.GetArticleByID(articleID)
	if err != nil {
		ah.logger.Println("Error getting article by ID:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if existingArticle == nil {
		http.NotFound(w, r)
		return
	}
	var updatedArticleRequest struct {
		Title       *string        `json:"title"`
		Description *string        `json:"description"`
		Image       *string        `json:"image"`
		AuthorID    *int           `json:"author_id"`
		Paraghraps  []store.Paraghraph `json:"paraghraps"`
	}
	err = json.NewDecoder(r.Body).Decode(&updatedArticleRequest)
	if err != nil {
		ah.logger.Println("error while decoding article:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	if updatedArticleRequest.Title != nil {
		existingArticle.Title = *updatedArticleRequest.Title
	}
	if updatedArticleRequest.Description != nil {
		existingArticle.Description = *updatedArticleRequest.Description
	}
	if updatedArticleRequest.Image != nil {
		existingArticle.Image = *updatedArticleRequest.Image
	}
	if updatedArticleRequest.AuthorID != nil {
		existingArticle.AuthorID = *updatedArticleRequest.AuthorID
	}
	if updatedArticleRequest.Paraghraps != nil {
		existingArticle.Paraghraps = updatedArticleRequest.Paraghraps
	}
	err = ah.articleStore.UpdateArticle(existingArticle)
	if err != nil {
		ah.logger.Println("Error updating article:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"article": existingArticle})
}

func (ah *ArticleHandler) HandleDeleteArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := utils.ReadIDParam(r)
	if err != nil {
		ah.logger.Println("Error reading article ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid article ID"})
		return
	}

	err = ah.articleStore.DeleteArticle(articleID)
	if err != nil {
		ah.logger.Println("Error deleting article:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}
