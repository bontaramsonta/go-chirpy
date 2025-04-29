package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bontaramsonta/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	// parse request body
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	// validations and defaults
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required", nil)
		return
	}
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > int(time.Hour.Seconds()) {
		params.ExpiresInSeconds = int(time.Hour.Seconds())
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

	// generate token
	token, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Duration(params.ExpiresInSeconds)*time.Second,
	)
	if err != nil {
		log.Println("Error generating token:", err)
		respondWithError(w, http.StatusInternalServerError, "Error generating token", err)
		return
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
