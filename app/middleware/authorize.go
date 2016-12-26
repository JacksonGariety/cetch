package middleware

import (
	"net/http"
	"fmt"
	"context"
	"github.com/dgrijalva/jwt-go"
	"os"

	"github.com/JacksonGariety/cetch/app/models"
)

var sessionKey = os.Getenv("session_key")
var sessionHash = os.Getenv("session_hash")

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Auth")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		// Return a Token using the cookie
		token, err := jwt.ParseWithClaims(cookie.Value, &models.Claims{}, func(token *jwt.Token) (interface{}, error){
			// Make sure token's signature wasn't changed
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected siging method")
			}
			return []byte(sessionHash), nil
		})
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), sessionKey, *claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
			return
		}
	})
}

func CurrentUser(r *http.Request) (models.Claims, bool) {
	claims, ok := r.Context().Value(sessionKey).(models.Claims)
	return claims, ok
}