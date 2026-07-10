package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/imatimba/http-server/internal/auth"
	"github.com/imatimba/http-server/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	bearerToken, err := auth.GetAuthToken(r.Header)
	if err != nil {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	userUUID, err := auth.ValidateJWT(bearerToken, cfg.secretKey)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	if !validateChirp(params.Body) {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Chirp contains blocked words"})
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userUUID,
	})

	if err != nil {
		http.Error(w, "Failed to create chirp", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve chirps", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		http.Error(w, "Failed to retrieve chirp", http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		http.Error(w, "Invalid chirp ID", http.StatusBadRequest)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpUUID)
	if err != nil {
		http.Error(w, "Failed to retrieve chirp", http.StatusNotFound)
		return
	}

	bearerToken, err := auth.GetAuthToken(r.Header)
	if err != nil {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	userUUID, err := auth.ValidateJWT(bearerToken, cfg.secretKey)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	if chirp.UserID != userUUID {
		http.Error(w, "You are not authorized to delete this chirp", http.StatusForbidden)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpUUID)
	if err != nil {
		http.Error(w, "Failed to delete chirp", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
