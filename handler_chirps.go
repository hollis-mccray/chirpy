package main

import (
	"encoding/json"
	"net/http"
	"sort"
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
		respondWithError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanText := contentFilter(params.Body)

	response, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanText,
		UserID: id,
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
		UserID:    response.UserID,
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
		UserID:    response.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *apiConfig) listAllChirps(w http.ResponseWriter, r *http.Request) {

	response, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	var chirps []Chirp

	for _, dbChirp := range response {
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		chirps = append(chirps, chirp)
	}

	order := r.URL.Query().Get("sort")

	if strings.ToLower(order) == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	} else {
		sort.Slice(chirps, func(i, j int) bool { return chirps[j].CreatedAt.After(chirps[i].CreatedAt) })
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
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

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtkey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if response.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Forbidden", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), response.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
