package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bontaramsonta/go-chirpy/internal/auth"
	"github.com/bontaramsonta/go-chirpy/internal/database"
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
	// validations and defaults
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

	// generate token
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		log.Println("Error generating token:", err)
		respondWithError(w, http.StatusInternalServerError, "Error generating token", err)
		return
	}

	// generate refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Println("Error generating refresh token:", err)
		respondWithError(w, http.StatusInternalServerError, "Error generating refresh token", err)
		return
	}

	// save refresh token
	if err := cfg.db.SaveRefreshToken(r.Context(), database.SaveRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * auth.RefreshTokenExpirationDays),
	}); err != nil {
		log.Println("Error saving refresh token:", err)
		respondWithError(w, http.StatusInternalServerError, "Error saving refresh token", err)
		return
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
