package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/imatimba/http-server/internal/database"
)

type userResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Email     string `json:"email"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func validateChirp(body string) bool {
	blockedWords := []string{"kerfuffle", "sharbert", "fornax"}
	bodyWords := strings.Split(body, " ")

	for _, word := range bodyWords {
		for _, blocked := range blockedWords {
			if strings.EqualFold(word, blocked) {
				return false
			}
		}
	}

	return true
}

func userToResponse(user database.User) userResponse {
	return userResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Email:     user.Email,
	}
}
