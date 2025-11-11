package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/makhammatovb/Articles/internal/store"
	"github.com/makhammatovb/Articles/internal/utils"
)

type ReviewHandler struct {
	reviewStore store.ReviewStore
	articleStore store.ArticleStore
	logger      *log.Logger
}

func NewReviewHandler(reviewStore store.ReviewStore, articleStore store.ArticleStore, logger *log.Logger) *ReviewHandler {
	return &ReviewHandler{
		reviewStore: reviewStore,
		articleStore: articleStore,
		logger:      logger,
	}
}

func (rh *ReviewHandler) HandleGetReviewByID(w http.ResponseWriter, r *http.Request) {

	reviewID, err := utils.ReadIDParam(r)
	if err != nil {
		rh.logger.Println("Error reading review ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid review ID"})
		return
	}
	review, err := rh.reviewStore.GetReviewByID(reviewID)
	if err != nil {
		rh.logger.Println("Error getting review by ID:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"review": review})
}

func (rh *ReviewHandler) HandleCreateReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	var review store.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		rh.logger.Println("Decoding error:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	review.AuthorID = userID
	articleExists, err := rh.articleStore.ArticleExists(review.ArticleID)
    if err != nil {
        rh.logger.Println("Error checking article existence:", err)
        utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
        return
    }
    if !articleExists {
        utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Article not found"})
        return
    }
	articleAuthorID, err := rh.articleStore.GetArticleAuthorID(review.ArticleID)
    if err != nil {
        rh.logger.Println("Error getting article author:", err)
        utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
        return
    }
    if articleAuthorID == userID {
        utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "Cannot review your own article"})
        return
    }
	existingReview, err := rh.reviewStore.GetReviewByUserAndArticle(userID, review.ArticleID)
    if err != nil {
        rh.logger.Println("Error checking existing review:", err)
        utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
        return
    }
    if existingReview != nil {
        utils.WriteJSON(w, http.StatusConflict, utils.Envelope{"error": "You have already reviewed this article"})
        return
    }
	createdReview, err := rh.reviewStore.CreateReview(&review)
	if err != nil {
		rh.logger.Println("Error creating review:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"review": createdReview})
}

func (rh *ReviewHandler) HandleUpdateReview(w http.ResponseWriter, r *http.Request) {
	reviewID, err := utils.ReadIDParam(r)
	if err != nil {
		rh.logger.Println("Error reading review ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid review ID"})
		return
	}
	existingReview, err := rh.reviewStore.GetReviewByID(reviewID)
	if err != nil {
		rh.logger.Println("Error getting review by ID:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if existingReview == nil {
		http.NotFound(w, r)
		return
	}
	var updatedReviewRequest struct {
		Stars *int    `json:"stars"`
		Note  *string `json:"note"`
	}
	err = json.NewDecoder(r.Body).Decode(&updatedReviewRequest)
	if err != nil {
		rh.logger.Println("error while decoding review:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	if updatedReviewRequest.Stars != nil {
		existingReview.Stars = *updatedReviewRequest.Stars
	}
	if updatedReviewRequest.Note != nil {
		existingReview.Note = updatedReviewRequest.Note
	}
	err = rh.reviewStore.UpdateReview(existingReview)
	if err != nil {
		rh.logger.Println("Error updating review:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"review": existingReview})
}

func (rh *ReviewHandler) HandleDeleteReview(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int64)
	if !ok {
		rh.logger.Println("User ID not found in context")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Authentication required"})
		return
	}

	reviewID, err := utils.ReadIDParam(r)
	if err != nil {
		rh.logger.Println("Error reading review ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid review ID"})
		return
	}
	
	existingReview, err := rh.reviewStore.GetReviewByID(reviewID)
	if err != nil {
		rh.logger.Println("Error getting review by ID:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if existingReview == nil {
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Review not found"})
		return
	}

	if existingReview.AuthorID != userID {
		rh.logger.Printf("User %d attempted to delete review %d owned by user %d", userID, reviewID, existingReview.AuthorID)
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "You can only delete your own reviews"})
		return
	}

	err = rh.reviewStore.DeleteReview(reviewID)
	if err != nil {
		rh.logger.Println("Error deleting review:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}
