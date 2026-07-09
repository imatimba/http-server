package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/imatimba/http-server/internal/auth"
	"github.com/imatimba/http-server/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	ExpiresInSecs := 3600
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
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

	token, err := auth.MakeJWT(user.ID, cfg.secretKey, time.Duration(ExpiresInSecs)*time.Second)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})

	if err != nil {
		http.Error(w, "Failed to create refresh token", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, userToResponse(user, token, refreshToken))
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetAuthToken(r.Header)

	if err != nil {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	dbRefreshToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil || dbRefreshToken.RevokedAt.Valid || dbRefreshToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.secretKey, 3600*time.Second)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetAuthToken(r.Header)

	if err != nil {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token:     refreshToken,
		RevokedAt: sql.NullTime{Valid: true, Time: time.Now()},
		UpdatedAt: time.Now(),
	})

	if err != nil {
		http.Error(w, "Failed to revoke refresh token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
