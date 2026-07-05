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

func cleanBodyString(body string) string {
	blockedWords := []string{"kerfuffle", "sharbert", "fornax"}
	bodyWords := strings.Split(body, " ")

	for _, blockedWord := range blockedWords {
		for i, bodyWord := range bodyWords {
			if strings.ToLower(bodyWord) == blockedWord {
				bodyWords[i] = "****"
			}
		}
	}

	return strings.Join(bodyWords, " ")
}
