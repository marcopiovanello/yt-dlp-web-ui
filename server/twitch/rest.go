package twitch

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
)

type addUserReq struct {
	User string `json:"user"`
}

func MonitorUserHandler(m *Monitor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req addUserReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		m.Add(req.User)

		if err := json.NewEncoder(w).Encode("ok"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetMonitoredUsers(m *Monitor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		it := m.GetMonitoredUsers()

		users := slices.Collect(it)
		if users == nil {
			users = make([]string, 0)
		}

		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeleteUser(m *Monitor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := chi.URLParam(r, "user")

		if user == "" {
			http.Error(w, "empty user", http.StatusBadRequest)
			return
		}

		m.DeleteUser(user)

		if err := json.NewEncoder(w).Encode("ok"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
