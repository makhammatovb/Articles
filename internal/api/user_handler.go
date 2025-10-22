package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/makhammatovb/Articles/internal/store"
	"github.com/makhammatovb/Articles/internal/utils"
)

// UserHandler struct to handle User-related requests for future use
type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

// HandleGetUserByID handles the GET request to retrieve a User by its ID.
// wh *UserHandler is the receiver
// *UserHandler means the method operates on a pointer to UserHandler
// w http.ResponseWriter is used to write the response back to the client
// r *http.Request represents the incoming HTTP request
func (uh *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {

	// retrieves the user ID from the URL parameters
	userID, err := utils.ReadIDParam(r)
	fmt.Println("USER ID:", userID)
	if err != nil {
		uh.logger.Println("Error reading user ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}
	user, err := uh.userStore.GetUserByID(userID)
	if err != nil {
		uh.logger.Println("Error getting user by ID:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

// HandleCreateUser handles the POST request to create a new user.
func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uh.logger.Println("Decoding error:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	createdUser, err := uh.userStore.CreateUser(&user)
	if err != nil {
		uh.logger.Println("Error creating user:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": createdUser})
}

func (uh *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Println("Error reading user ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}
	existingUser, err := uh.userStore.GetUserByID(userID)
	if err != nil {
		uh.logger.Println("Error getting user by ID:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if existingUser == nil {
		http.NotFound(w, r)
		return
	}
	var updatedUserRequest struct {
		FirstName    *string `json:"firstname"`
		LastName     *string `json:"lastname"`
		Email        *string `json:"email"`
		PasswordHash *string `json:"password"`
	}
	err = json.NewDecoder(r.Body).Decode(&updatedUserRequest)
	if err != nil {
		uh.logger.Println("error while decoding user:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	if updatedUserRequest.FirstName != nil {
		existingUser.FirstName = *updatedUserRequest.FirstName
	}
	if updatedUserRequest.LastName != nil {
		existingUser.LastName = *updatedUserRequest.LastName
	}
	if updatedUserRequest.Email != nil {
		existingUser.Email = *updatedUserRequest.Email
	}
	if updatedUserRequest.PasswordHash != nil {
		existingUser.PasswordHash = *updatedUserRequest.PasswordHash
	}
	err = uh.userStore.UpdateUser(existingUser)
	if err != nil {
		uh.logger.Println("Error updating user:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": existingUser})
}

func (uh *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadIDParam(r)
	if err != nil {
		uh.logger.Println("Error reading user ID:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	err = uh.userStore.DeleteUser(userID)
	if err != nil {
		uh.logger.Println("Error deleting user:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusNoContent, nil)
}
