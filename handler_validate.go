package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type response struct {
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}
	
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Body: contentFilter(params.Body),
	})
}

func contentFilter(s string) string {
	words := strings.Split(s, " ")
	bad_words := []string {
		"kerfuffle", "sharbert", "fornax",
	}

	for i, word := range(words) {
		lowerCase := strings.ToLower(word)
		for _, bad_word := range bad_words{
			if lowerCase == bad_word {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}