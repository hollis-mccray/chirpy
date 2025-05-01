package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hollis-mccray/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, err = cfg.db.UpdateToRed(r.Context(), params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error procssing request", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
