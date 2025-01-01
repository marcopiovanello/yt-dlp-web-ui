package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive/domain"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/openid"

	middlewares "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/middleware"
)

type Handler struct {
	service domain.Service
}

func New(service domain.Service) domain.RestHandler {
	return &Handler{
		service: service,
	}
}

// List implements domain.RestHandler.
func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		var (
			startRowIdParam = r.URL.Query().Get("id")
			LimitParam      = r.URL.Query().Get("limit")
		)

		startRowId, err := strconv.Atoi(startRowIdParam)
		if err != nil {
			startRowId = 0
		}

		limit, err := strconv.Atoi(LimitParam)
		if err != nil {
			limit = 50
		}

		res, err := h.service.List(r.Context(), startRowId, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Archive implements domain.RestHandler.
func (h *Handler) Archive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		var req domain.ArchiveEntry

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := h.service.Archive(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode("ok")
	}
}

// HardDelete implements domain.RestHandler.
func (h *Handler) HardDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParam(r, "id")

		res, err := h.service.HardDelete(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// SoftDelete implements domain.RestHandler.
func (h *Handler) SoftDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParam(r, "id")

		res, err := h.service.SoftDelete(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetCursor implements domain.RestHandler.
func (h *Handler) GetCursor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParam(r, "id")

		cursorId, err := h.service.GetCursor(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(cursorId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ApplyRouter implements domain.RestHandler.
func (h *Handler) ApplyRouter() func(chi.Router) {
	return func(r chi.Router) {
		if config.Instance().RequireAuth {
			r.Use(middlewares.Authenticated)
		}
		if config.Instance().UseOpenId {
			r.Use(openid.Middleware)
		}

		r.Get("/", h.List())
		r.Get("/cursor/{id}", h.GetCursor())
		r.Post("/", h.Archive())
		r.Delete("/soft/{id}", h.SoftDelete())
		r.Delete("/hard/{id}", h.HardDelete())
	}
}
