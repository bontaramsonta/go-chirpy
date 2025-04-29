package main

import (
	"net/http"

	"github.com/bontaramsonta/go-chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersRefresh(w http.ResponseWriter, r *http.Request) {
	// get userID from context
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	// generate access token
	accessToken, err := auth.MakeJWT(userID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	type response struct {
		AccessToken string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{AccessToken: accessToken})
}
