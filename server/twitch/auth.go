package twitch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const authURL = "https://id.twitch.tv/oauth2/token"

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AccessToken struct {
	Token  string
	Expiry time.Time
}

type AuthenticationManager struct {
	clientId     string
	clientSecret string
	accesToken   *AccessToken
}

func NewAuthenticationManager(clientId, clientSecret string) *AuthenticationManager {
	return &AuthenticationManager{
		clientId:     clientId,
		clientSecret: clientSecret,
		accesToken:   &AccessToken{},
	}
}

func (a *AuthenticationManager) GetAccessToken() (*AccessToken, error) {
	if a.accesToken != nil && a.accesToken.Expiry.After(time.Now()) {
		return a.accesToken, nil
	}

	data := url.Values{}
	data.Set("client_id", a.clientId)
	data.Set("client_secret", a.clientSecret)
	data.Set("grant_type", "client_credentials")

	resp, err := http.PostForm(authURL, data)
	if err != nil {
		return nil, fmt.Errorf("errore richiesta token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status non OK: %s", resp.Status)
	}

	var auth AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return nil, fmt.Errorf("errore decoding JSON: %w", err)
	}

	token := &AccessToken{
		Token:  auth.AccessToken,
		Expiry: time.Now().Add(time.Duration(auth.ExpiresIn) * time.Second),
	}

	a.accesToken = token

	return token, nil
}

func (a *AuthenticationManager) GetClientId() string {
	return a.clientId
}
