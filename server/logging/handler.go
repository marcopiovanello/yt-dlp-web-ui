package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/marcopeocchi/yt-dlp-web-ui/server/config"
	middlewares "github.com/marcopeocchi/yt-dlp-web-ui/server/middleware"
	"github.com/marcopeocchi/yt-dlp-web-ui/server/openid"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1000,
	WriteBufferSize: 1000,
}

func webSocket(logger *ObservableLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logs := logger.Observe(r.Context())

		for {
			select {
			case <-r.Context().Done():
				return
			case msg := <-logs:
				c.WriteJSON(msg)
			}
		}
	}
}

func sse(logger *ObservableLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		logs := logger.Observe(r.Context())

		for {
			select {
			case <-r.Context().Done():
				slog.Info("detaching from logger")
				return
			case msg, ok := <-logs:
				if !ok {
					http.Error(w, "closed logs channel", http.StatusInternalServerError)
					return
				}

				var (
					b  bytes.Buffer
					sb strings.Builder
				)

				if err := json.NewEncoder(&b).Encode(msg); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				sb.WriteString("event: log\n")
				sb.WriteString("data: ")
				sb.WriteString(b.String())
				sb.WriteRune('\n')
				sb.WriteRune('\n')

				fmt.Fprint(w, sb.String())

				flusher.Flush()
			}
		}
	}
}

func ApplyRouter(logger *ObservableLogger) func(chi.Router) {
	return func(r chi.Router) {
		if config.Instance().RequireAuth {
			r.Use(middlewares.Authenticated)
		}
		if config.Instance().UseOpenId {
			r.Use(openid.Middleware)
		}
		r.Get("/ws", webSocket(logger))
		r.Get("/sse", sse(logger))
	}
}
