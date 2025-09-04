package middlewares

import (
	"net/http"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/openid"
)

func ApplyAuthenticationByConfig(next http.Handler) http.Handler {
	handler := next

	if config.Instance().Authentication.RequireAuth {
		handler = Authenticated(handler)
	}
	if config.Instance().OpenId.UseOpenId {
		handler = openid.Middleware(handler)
	}

	return handler
}
