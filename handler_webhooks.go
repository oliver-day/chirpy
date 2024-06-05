package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/oday/chirpy/internal/auth"
	"github.com/oday/chirpy/internal/database"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	polkaApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to find API key")
		return
	}
	if polkaApiKey != cfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.DB.UpgradeChirpyRed(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrorUserDoesNotExist) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
