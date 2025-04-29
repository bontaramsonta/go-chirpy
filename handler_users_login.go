package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bontaramsonta/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	// parse request body
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// validations
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required", nil)
		return
	}
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}

	authenticationErrResponse := func(err error) {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
	}

	// authenticate user
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Println("Error getting user:", err)
		authenticationErrResponse(err)
		return
	}
	if err := auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		// log.Println("Error checking password:", err)
		authenticationErrResponse(err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
