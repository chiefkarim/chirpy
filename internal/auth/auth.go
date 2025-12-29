package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	isTrue, err := argon2id.ComparePasswordAndHash(password, hash)
	return isTrue, err
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))
	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Printf("error parsing jwt token %v", err)
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		if userid, ok := uuid.Parse(claims.Subject); ok == nil {
			return userid, nil
		} else {
			log.Print("Unkonw claims type")
			err = fmt.Errorf("Unkown claims type jwt.RegisteredClaims")
		}
	}

	return uuid.UUID{}, err
}

func GetBearerToken(header http.Header) (string, error) {
	res := header.Get("Authorization")
	if res == "" {
		return "", errors.New("authorization header not present")
	}
	authHeader := strings.Split(res, " ")
	if len(authHeader) != 2 {
		return "", errors.New("wrong formated authorization header")
	}

	return authHeader[1], nil
}

func MakeRefreshToken() (string, error) {
	var random []byte
	_, err := rand.Read(random)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(random), nil
}
