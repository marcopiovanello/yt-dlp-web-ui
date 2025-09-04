package openid

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"golang.org/x/oauth2"
)

var (
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
)

func Configure() {
	if !config.Instance().OpenId.UseOpenId {
		return
	}

	provider, err := oidc.NewProvider(
		context.Background(),
		config.Instance().OpenId.ProviderURL,
	)
	if err != nil {
		panic(err)
	}

	oauth2Config = oauth2.Config{
		ClientID:     config.Instance().OpenId.ClientId,
		ClientSecret: config.Instance().OpenId.ClientSecret,
		RedirectURL:  config.Instance().OpenId.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	verifier = provider.Verifier(&oidc.Config{
		ClientID: config.Instance().OpenId.ClientId,
	})
}
