package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	TokenIssuer                = "chirpy"
	UserIDKey                  = "userID"
	RefreshTokenKey            = "refreshToken"
	AccessTokenExpiration      = time.Hour
	RefreshTokenExpirationDays = 60
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    TokenIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
		Subject:   userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func MakeRefreshToken() (string, error) {
	// refresh token is a random 256-bit (32-byte) hex-encoded string
	randBytes := make([]byte, 32)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}
	refreshToken := hex.EncodeToString(randBytes)
	return refreshToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// parse token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)

	invalidTokenErr := fmt.Errorf("invalid token")
	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err)
		return uuid.Nil, invalidTokenErr
	}

	// parse subject
	sub, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("Failed to get subject from token: %v", err)
		return uuid.Nil, invalidTokenErr
	}

	// parse UUID from subject
	userID, err := uuid.Parse(sub)
	if err != nil {
		log.Printf("Failed to parse subject as UUID: %v", err)
		return uuid.Nil, invalidTokenErr
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", fmt.Errorf("invalid authorization header format should be: ApiKey <API_KEY>")
	}

	return parts[1], nil
}
