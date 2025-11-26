package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/makhammatovb/Articles/internal/store"
	"github.com/makhammatovb/Articles/internal/tokens"
	"github.com/makhammatovb/Articles/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest

	h.logger.Println(r.Body)

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Println("error while decoding token:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		h.logger.Println("error while getting user:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Println("error while comparing passwords:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	if !passwordsDoMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(int64(user.ID), 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Println("error while creating token:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"token": token})
}

func (h *TokenHandler) GenerateResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Println("error while decoding password reset request:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	if req.Email == "" {
		h.logger.Println("email is required for password reset")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Email is required"})
		return
	}
	
	user, err := h.userStore.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.Println("Error reading user email:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user email"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(int64(user.ID), 60*time.Minute, "reset-password")
	if err != nil {
		h.logger.Println("error while creating reset token:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"reset_token": token})
}

func (h *TokenHandler) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	token, err := utils.ReadTokenParam(r)
    if token == "" {
        utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Reset token is required"})
        return
    }
	if err != nil {
		h.logger.Println("Error reading reset token:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid reset token"})
		return
	}
	var req struct {
		NewPassword string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Println("error while decoding reset password request:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	if req.NewPassword == "" || req.ConfirmPassword == "" {
        utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "New password and confirmation are required"})
        return
    }
	if req.NewPassword != req.ConfirmPassword {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Passwords do not match"})
		return
	}
	
	tokenData, err := h.tokenStore.GetToken(token, "reset-password")
	if err != nil || tokenData == nil  || tokenData.Expiry.Before(time.Now()) {
		h.logger.Println("Error retrieving reset token:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid or expired reset token"})
		return
	}

	user, err := h.userStore.GetUserWithPasswordByID(tokenData.UserID)
    if err != nil {
        h.logger.Println("Error getting user by ID:", err)
        utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
        return
    }

	err = user.PasswordHash.Set(req.NewPassword)
    if err != nil {
        h.logger.Println("Error hashing new password:", err)
        utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
        return
    }

    err = h.userStore.UpdateUser(user)
    if err != nil {
        h.logger.Println("Error updating password:", err)
        utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
        return
    }

	err = h.tokenStore.DeleteToken(token, "reset-password")
    if err != nil {
        h.logger.Println("Error deleting used token:", err)
    }
    utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Password reset successfully"})
}
