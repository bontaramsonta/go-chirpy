package main

import (
	"context"
	"net/http"

	"github.com/bontaramsonta/go-chirpy/internal/auth"
)

func (cfg *apiConfig) middlewareisAuthed(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
			return
		}

		userID, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *apiConfig) middlewareCheckRefreshToken(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
			return
		}

		// get user from refresh token
		userID, err := cfg.db.GetUserIdFromValidRefreshToken(r.Context(), refreshToken)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
			return
		}

		ctx := context.WithValue(r.Context(), auth.UserIDKey, userID)
		ctx = context.WithValue(ctx, auth.RefreshTokenKey, refreshToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
