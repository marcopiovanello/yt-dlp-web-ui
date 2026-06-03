package rpc

import (
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSockets JSON-RPC handler
func WebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer c.Close()

	// notify client that conn is open and ok
	c.WriteJSON(struct{ Status string }{Status: "connected"})

	for {
		mtype, reader, err := c.NextReader()
		if err != nil {
			break
		}

		res := newRequest(reader).Call()

		writer, err := c.NextWriter(mtype)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		io.Copy(writer, res)
	}
}

// HTTP-POST JSON-RPC handler
func Post(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	res := newRequest(r.Body).Call()
	_, err := io.Copy(w, res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
