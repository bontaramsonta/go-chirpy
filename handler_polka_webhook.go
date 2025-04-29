package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const (
	EventUserUpgraded = "user.upgraded"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// parse body
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
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
