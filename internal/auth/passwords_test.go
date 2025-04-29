package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	err = CheckPasswordHash(hashedPassword, password)
	if err != nil {
		t.Errorf("Error checking password: %v", err)
	}
}

func TestCheckPassword(t *testing.T) {
	password := "password123"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	err = CheckPasswordHash(hashedPassword, "wrongpassword")
	if err == nil {
		t.Errorf("Expected error for wrong password")
	}
}

func TestCheckPasswordWithEmptyPassword(t *testing.T) {
	hashedPassword, err := HashPassword("")
	if err != nil {
		t.Errorf("Error hashing empty password: %v", err)
	}

	err = CheckPasswordHash(hashedPassword, "")
	if err != nil {
		t.Errorf("Error checking empty password: %v", err)
	}
}

func TestCheckPasswordWithInvalidHash(t *testing.T) {
	password := "password123"
	_, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	err = CheckPasswordHash("invalidhash", password)
	if err == nil {
		t.Errorf("Expected error for invalid hash")
	}
}
