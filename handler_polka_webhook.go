package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/bontaramsonta/go-chirpy/internal/auth"
	"github.com/google/uuid"
)

const (
	EventUserUpgraded = "user.upgraded"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// check API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", errors.New("polka API key mismatch"))
		return
	}

	// parse body
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	params := parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	switch params.Event {
	case EventUserUpgraded:
		_, err := cfg.db.UpgradeUser(r.Context(), params.Data.UserID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				log.Print("polka webhook invalid userid")
				break
			}
			respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user", err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
