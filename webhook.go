package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/imatimba/http-server/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, "Missing or invalid API key", http.StatusUnauthorized)
		return
	}
	if apiKey != cfg.polkaKey {
		http.Error(w, "Invalid API key", http.StatusUnauthorized)
		return
	}

	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if params.Event != "user.upgraded" {
		http.Error(w, "Invalid event type", http.StatusNoContent)
		return
	}

	userUUID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := cfg.db.LookupUserByID(r.Context(), userUUID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to upgrade user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
