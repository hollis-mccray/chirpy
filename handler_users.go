package main

import (
	"encoding/json"
	"net/http"

	"github.com/hollis-mccray/chirpy/internal/auth"
	"github.com/hollis-mccray/chirpy/internal/database"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	pwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
	}

	user := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: pwd,
	}

	response, err := cfg.db.CreateUser(r.Context(), user)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:          response.ID,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
		Email:       response.Email,
		IsChirpyRed: response.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

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

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	updateParams := database.UpdatePasswordParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	response, err := cfg.db.UpdatePassword(r.Context(), updateParams)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          response.ID,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
		Email:       response.Email,
		IsChirpyRed: response.IsChirpyRed,
	})
}
