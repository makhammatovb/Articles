package utils

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// JSON format representor, keys are strings, values can be any type
type Envelope map[string]interface{}

// WriteJSON sends JSON response to the client
func WriteJSON(w http.ResponseWriter, status int, data Envelope) {
	// MarshalIndent formats the JSON with indentation for better readability
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js = append(js, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

// ReadIDParam reads the "id" parameter from the URL and converts it to int64
func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, http.ErrNoLocation
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, http.ErrNoLocation
	}
	return id, nil
}
