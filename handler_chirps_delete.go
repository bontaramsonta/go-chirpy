package main

import (
	"net/http"
	"strconv"

	"github.com/bontaramsonta/go-chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// get userID from context
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	// get chirpID from path
	chirpId := r.PathValue("chirpID")
	if chirpId == "" {
		respondWithError(w, http.StatusBadRequest, "Chirp ID is required", nil)
		return
	}

	id, err := strconv.Atoi(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// get chirp from database
	dbChirp, err := cfg.db.GetChirpByID(r.Context(), int32(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	// check if user is the owner of the chirp
	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not the owner of this chirp", nil)
		return
	}

	// delete chirp from database
	err = cfg.db.DeleteChirp(r.Context(), int32(id))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
