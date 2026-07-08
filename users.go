package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/imatimba/http-server/internal/auth"
	"github.com/imatimba/http-server/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:        params.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, userToResponse(user, ""))
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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password      string `json:"password"`
		Email         string `json:"email"`
		ExpiresInSecs int64  `json:"expires_in_seconds"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if params.ExpiresInSecs <= 0 || params.ExpiresInSecs > 3600 {
		params.ExpiresInSecs = 3600 // default to 1 hour if not provided or invalid
	}

	user, err := cfg.db.LookupUserByEmail(r.Context(), params.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	isValid, err := auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil || !isValid {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(params.ExpiresInSecs)*time.Second)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, userToResponse(user, token))
}
