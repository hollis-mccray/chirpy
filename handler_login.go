package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hollis-mccray/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Expires  int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	if params.Expires == 0 || params.Expires > 3600 {
		params.Expires = 3600
	}

	response, err := cfg.db.UserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", params.Expires))
	token, err := auth.MakeJWT(response.ID, cfg.jwtkey, duration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error handling request", err)
		return
	}

	err = auth.CheckPasswordHash(response.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
	} else {
		respondWithJSON(w, http.StatusOK, User{
			ID:        response.ID,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
			Email:     response.Email,
			Token:     token,
		})
	}
}
