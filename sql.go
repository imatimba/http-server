package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		http.Error(w, "Forbidden, not running in dev mode", http.StatusForbidden)
		return
	}

	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to delete users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset OK"))
}
