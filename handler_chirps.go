package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hollis-mccray/chirpy/internal/auth"
	"github.com/hollis-mccray/chirpy/internal/database"
)

func (cfg *apiConfig) handlerNewChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	id, err := auth.ValidateJWT(token, cfg.jwtkey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid auth token", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanText := contentFilter(params.Body)

	response, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanText,
		UserID: uuid.NullUUID{
			UUID:  id,
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        response.ID,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
		Body:      response.Body,
		UserID: uuid.NullUUID{
			UUID:  id,
			Valid: true,
		},
	})
}

func contentFilter(s string) string {
	words := strings.Split(s, " ")
	bad_words := []string{
		"kerfuffle", "sharbert", "fornax",
	}

	for i, word := range words {
		lowerCase := strings.ToLower(word)
		for _, bad_word := range bad_words {
			if lowerCase == bad_word {
				words[i] = "****"
			}
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	response, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	chirp := Chirp{
		ID:        response.ID,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
		Body:      response.Body,
		UserID: uuid.NullUUID{
			UUID:  response.UserID.UUID,
			Valid: true,
		},
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) listAllChirps(w http.ResponseWriter, r *http.Request) {
	response, err := cfg.db.ListAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	var chirps []Chirp

	for _, item := range response {
		chirp := Chirp{
			ID:        item.ID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			Body:      item.Body,
			UserID: uuid.NullUUID{
				UUID:  item.UserID.UUID,
				Valid: true,
			},
		}
		chirps = append(chirps, chirp)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
