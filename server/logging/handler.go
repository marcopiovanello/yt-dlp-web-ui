package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	middlewares "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/middleware"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/openid"
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

				var b bytes.Buffer

				b.WriteString("event: log\n")
				b.WriteString("data: ")

				if err := json.NewEncoder(&b).Encode(msg); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				b.WriteRune('\n')
				b.WriteRune('\n')

				io.Copy(w, &b)

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
