package main

import (
	"encoding/json"
	"net/http"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(params.Body) > 140 {
		respondWithJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Chirp is too long",
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"cleaned_body": cleanBodyString(params.Body),
	})
}
