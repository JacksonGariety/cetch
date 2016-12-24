package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

func ClaimsCreate(username string) (string, time.Time) {

	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Hour * 1)

	claims := Claims {
		username,
		jwt.StandardClaims {
			ExpiresAt: expireToken,
			Issuer:    "localhost:8080",
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := token.SignedString([]byte("secret"))

	return signedToken, expireCookie
}
