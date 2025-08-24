package twitch

import (
	"encoding/json"
	"net/http"
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
