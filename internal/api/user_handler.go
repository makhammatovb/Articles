package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/makhammatovb/Articles/internal/store"
	"github.com/makhammatovb/Articles/internal/utils"
)

type registerUserRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

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

func (uh *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Email == "" {
		return errors.New("missing required fields")
	}

	if len(req.Email) > 255 {
		return errors.New("email is too long")
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	if req.Password == "" {
		return errors.New("missing required fields")
	}

	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Println("error while decoding user:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	err = uh.validateRegisterRequest(&req)
	if err != nil {
		uh.logger.Println("error while validating user:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Email: req.Email,
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Println("error while hashing password:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Println("error while creating user:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})

}

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
