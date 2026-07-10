package main

import (
	"encoding/json"
	"net/http"
	"sort"

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
	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		userUUID, err := uuid.Parse(authorID)
		if err != nil {
			http.Error(w, "Invalid author ID", http.StatusBadRequest)
			return
		}

		chirps, err := cfg.db.GetChirpsByUserID(r.Context(), userUUID)
		if err != nil {
			http.Error(w, "Failed to retrieve chirps", http.StatusInternalServerError)
			return
		}

		respondWithJSON(w, http.StatusOK, chirps)
		return
	}

	sortBy := r.URL.Query().Get("sort")

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve chirps", http.StatusInternalServerError)
		return
	}

	switch sortBy {
	case "desc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	case "asc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
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
