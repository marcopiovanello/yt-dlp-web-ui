package middlewares

import (
	"net/http"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/openid"
)

func ApplyAuthenticationByConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.Instance().RequireAuth {
			Authenticated(next)
		}
		if config.Instance().UseOpenId {
			openid.Middleware(next)
		}
		next.ServeHTTP(w, r)
	})
}
