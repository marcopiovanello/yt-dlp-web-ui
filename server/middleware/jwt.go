package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func validateToken(tokenValue string) error {
	token, err := jwt.Parse(tokenValue, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiresAt, err := time.Parse(time.RFC3339, claims["expiresAt"].(string))
		if err != nil {
			return err
		}

		if time.Now().After(expiresAt) {
			return errors.New("token expired")
		}
	} else {
		return errors.New("invalid token")
	}

	return nil
}

// Authentication does NOT use http-Only cookies since there's not risk for XSS
// By exposing the server through https it's completely safe to use httpheaders

func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Authentication")
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		if err := validateToken(token); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
