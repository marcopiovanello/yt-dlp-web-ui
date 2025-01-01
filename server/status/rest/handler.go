package rest

import (
	"encoding/json"
	"net/http"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/status/domain"
)

type RestHandler struct {
	service domain.Service
}

// Status implements domain.RestHandler.
func (h *RestHandler) Status() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json")

		status, err := h.service.Status(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func New(service domain.Service) domain.RestHandler {
	return &RestHandler{
		service: service,
	}
}
