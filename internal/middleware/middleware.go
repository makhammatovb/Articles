package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/makhammatovb/Articles/internal/store"
	"github.com/makhammatovb/Articles/internal/tokens"
	"github.com/makhammatovb/Articles/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) (*store.User, error) {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
        return nil, fmt.Errorf("user not found in context")
    }
    return user, nil
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid auth header format"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			fmt.Println("Error retrieving user by token:", err)
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid token"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetUser(r)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
			return
		}
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "you must be authenticated to access this resource"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
