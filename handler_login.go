package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hollis-mccray/chirpy/internal/auth"
	"github.com/hollis-mccray/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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

	response, err := cfg.db.UserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(response.ID, cfg.jwtkey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error handling request", err)
		return
	}

	refresh, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error handling request", err)
		return
	}

	now := time.Now()
	expires := now.Add(time.Hour * 1440)
	refreshParams := database.CreateRefreshTokenParams{
		Token:     refresh,
		CreatedAt: now,
		UpdatedAt: now,
		UserID: uuid.NullUUID{
			UUID:  response.ID,
			Valid: true,
		},
		ExpiresAt: expires,
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error handling request", err)
		return
	}

	err = auth.CheckPasswordHash(response.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	} else {
		respondWithJSON(w, http.StatusOK, User{
			ID:           response.ID,
			CreatedAt:    response.CreatedAt,
			UpdatedAt:    response.UpdatedAt,
			Email:        response.Email,
			Token:        token,
			RefreshToken: refresh,
		})
	}
}
