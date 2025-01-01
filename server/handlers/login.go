package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

const TOKEN_COOKIE_NAME = "jwt-yt-dlp-webui"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		username = config.Instance().Username
		password = config.Instance().Password
	)

	if username != req.Username || password != req.Password {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}

	expiresAt := time.Now().Add(time.Hour * 24 * 30)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiresAt": expiresAt,
		"username":  req.Username,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(tokenString); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     TOKEN_COOKIE_NAME,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now(),
		Value:    "",
		Path:     "/",
	}

	http.SetCookie(w, cookie)
}
