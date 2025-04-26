package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerNewUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	response, err := cfg.db.CreateUser(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User {
		ID: 		response.ID,
		CreatedAt: 	response.CreatedAt,
		UpdatedAt:	response.UpdatedAt,
		Email:		response.Email,
	})
}