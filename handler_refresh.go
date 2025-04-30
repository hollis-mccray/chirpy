package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/hollis-mccray/chirpy/internal/auth"
	"github.com/hollis-mccray/chirpy/internal/database"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	refreshToken, err := cfg.db.GetRefreshTokenFromToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	} else if refreshToken.RevokedAt.Valid && time.Now().After(refreshToken.RevokedAt.Time) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtkey, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	type accessResponse struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, accessResponse{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), database.RevokeTokenParams{
		Token: token,
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error handling request", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, nil)
}
