package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "mySecretPassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hashedPassword == "" {
		t.Error("Hashed password should not be empty")
	}

	// Test if the hashed password can be verified
	match, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed: %v", err)
	}
	if !match {
		t.Error("Password should match the hashed password")
	}

	// Test with a wrong password
	wrongPassword := "wrongPassword"
	match, err = CheckPasswordHash(wrongPassword, hashedPassword)
	if err != nil {
		t.Fatalf("CheckPasswordHash with wrong password failed: %v", err)
	}
	if match {
		t.Error("Wrong password should not match the hashed password")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "supersecretjwtkey"
	expiresIn := time.Minute * 5

	// Test MakeJWT
	jwtToken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	if jwtToken == "" {
		t.Error("JWT token should not be empty")
	}

	// Test ValidateJWT with a valid token
	validatedUserID, err := ValidateJWT(jwtToken, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed for valid token: %v", err)
	}
	if validatedUserID != userID {
		t.Errorf("Validated UserID (%s) does not match original UserID (%s)", validatedUserID, userID)
	}

	// Test ValidateJWT with a wrong secret
	wrongSecret := "wrongsecret"
	_, err = ValidateJWT(jwtToken, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT should fail with wrong secret")
	}

	// Test ValidateJWT with an expired token (by making a token with a very short expiry)
	expiredExpiresIn := time.Millisecond * 1
	expiredJwtToken, err := MakeJWT(userID, tokenSecret, expiredExpiresIn)
	if err != nil {
		t.Fatalf("MakeJWT for expired token failed: %v", err)
	}
	time.Sleep(expiredExpiresIn * 2) // Wait for the token to expire

	_, err = ValidateJWT(expiredJwtToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT should fail for an expired token")
	}
}
