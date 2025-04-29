package main

import (
	"net/http"

	"github.com/bontaramsonta/go-chirpy/internal/auth"
	"github.com/bontaramsonta/go-chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersRevoke(w http.ResponseWriter, r *http.Request) {
	// get userID, refreshToken from context
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	refreshToken := r.Context().Value(auth.RefreshTokenKey).(string)

	// revoke refresh token
	err := cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		Token:  refreshToken,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
