package auth

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

	match, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Fatalf("CheckPasswordHash failed: %v", err)
	}
	if !match {
		t.Error("Password should match the hashed password")
	}

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

	jwtToken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	if jwtToken == "" {
		t.Error("JWT token should not be empty")
	}

	validatedUserID, err := ValidateJWT(jwtToken, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT failed for valid token: %v", err)
	}
	if validatedUserID != userID {
		t.Errorf("Validated UserID (%s) does not match original UserID (%s)", validatedUserID, userID)
	}

	wrongSecret := "wrongsecret"
	_, err = ValidateJWT(jwtToken, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT should fail with wrong secret")
	}

	expiredExpiresIn := time.Millisecond * 1
	expiredJwtToken, err := MakeJWT(userID, tokenSecret, expiredExpiresIn)
	if err != nil {
		t.Fatalf("MakeJWT for expired token failed: %v", err)
	}
	time.Sleep(expiredExpiresIn * 2)

	_, err = ValidateJWT(expiredJwtToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT should fail for an expired token")
	}
}

type TestCase struct {
	input  http.Header
	output string
	Error  error
}

func TestGetBearerToken(t *testing.T) {
	cases := map[string]TestCase{
		"clean beearer token": {
			input: func() http.Header {
				h := http.Header{}
				h.Add("Authorization", "Bearer asdfasdfasdfd")
				return h
			}(),
			output: "asdfasdfasdfd",
			Error:  nil,
		},
		"no authorization header": {
			input: func() http.Header {
				h := http.Header{}
				return h
			}(),
			output: "",
			Error:  errors.New("authorization header not present"),
		},
		"authorization header wrong format": {
			input: func() http.Header {
				h := http.Header{}
				h.Add("Authorization", "asdfasdfasdf")
				return h
			}(),
			Error:  errors.New("wrong formated authorization header"),
			output: "",
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			wantValue := test.output
			wantError := test.Error

			haveValue, haveError := GetBearerToken(test.input)
			switch {
			case wantError == nil && haveError != nil:
				t.Error(cmp.Diff(wantError, haveError.Error()))
			case wantError != nil && haveError == nil:
				t.Error(cmp.Diff(wantError.Error(), haveError))
			case wantError != nil && haveError != nil:
				cmp.Diff(wantError.Error(), haveError.Error())
			}

			valueDiff := cmp.Diff(haveValue, wantValue)

			if valueDiff != "" {
				t.Error(valueDiff)
			}
		})
	}
}
