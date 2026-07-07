package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
