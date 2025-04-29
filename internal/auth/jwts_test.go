package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// happy‐path: we can create a JWT and immediately validate it
func TestMakeAndValidateJWT(t *testing.T) {
	secret := "super-secret-key"
	userID := uuid.New()
	// give it a little bit of life
	tokenString, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("makeJWT returned error: %v", err)
	}

	gotID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned unexpected error: %v", err)
	}
	if gotID != userID {
		t.Errorf("ValidateJWT returned userID %q, want %q", gotID, userID)
	}
}

// expired‐token: even though we just generated it, it should be treated as expired
func TestValidateExpiredJWT(t *testing.T) {
	secret := "another-secret"
	userID := uuid.New()
	// expires in the past
	tokenString, err := MakeJWT(userID, secret, -time.Second)
	if err != nil {
		t.Fatalf("makeJWT returned error: %v", err)
	}

	gotID, err := ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatal("ValidateJWT did not return error for expired token")
	}
	if gotID != uuid.Nil {
		t.Errorf("ValidateJWT returned userID %q for expired token; want Nil", gotID)
	}
}

// wrong‐secret: token was not signed with the secret we pass to ValidateJWT
func TestValidateJWTWrongSecret(t *testing.T) {
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"
	userID := uuid.New()

	tokenString, err := MakeJWT(userID, correctSecret, time.Minute)
	if err != nil {
		t.Fatalf("makeJWT returned error: %v", err)
	}

	gotID, err := ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Fatal("ValidateJWT did not return error for wrong secret")
	}
	if gotID != uuid.Nil {
		t.Errorf("ValidateJWT returned userID %q for wrong-secret case; want Nil", gotID)
	}
}

// malformed token: totally not a JWT
func TestValidateMalformedToken(t *testing.T) {
	secret := "whatever"
	_, err := ValidateJWT("this-is-not-a-jwt", secret)
	if err == nil {
		t.Fatal("ValidateJWT did not return error for malformed token")
	}
}

// tampered token: flip a character in the signature part
func TestValidateTamperedToken(t *testing.T) {
	secret := "tamper-secret"
	userID := uuid.New()
	tokenString, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("makeJWT returned error: %v", err)
	}

	// JWTs are three base64‐URL parts separated by dots.
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		t.Fatalf("expected token to have 3 parts, got %d", len(parts))
	}
	// invert one bit in the signature
	sig := []rune(parts[2])
	sig[0] = sig[0] ^ 1
	parts[2] = string(sig)
	tampered := strings.Join(parts, ".")

	gotID, err := ValidateJWT(tampered, secret)
	if err == nil {
		t.Fatal("ValidateJWT did not return error on tampered token")
	}
	if gotID != uuid.Nil {
		t.Errorf("ValidateJWT returned userID %q for tampered token; want Nil", gotID)
	}
}
