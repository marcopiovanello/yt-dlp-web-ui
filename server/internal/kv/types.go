package kv

import "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"

// struct representing the current status of the memoryDB
// used for serializaton/persistence reasons
type Session struct {
	Processes []internal.ProcessSnapshot `json:"processes"`
}
