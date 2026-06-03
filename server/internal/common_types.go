package internal

import (
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"
)

// Used to unmarshall yt-dlp progress
type ProgressTemplate struct {
	Percentage string  `json:"percentage"`
	Speed      float64 `json:"speed"`
	Size       string  `json:"size"`
	Eta        float64 `json:"eta"`
}

type PostprocessTemplate struct {
	FilePath string `json:"filepath"`
}

// Defines where and how the download needs to be saved
type DownloadOutput struct {
	Path          string
	Filename      string
	SavedFilePath string `json:"savedFilePath"`
}

// Progress for the Running call
type DownloadProgress struct {
	Status     int     `json:"process_status"`
	Percentage string  `json:"percentage"`
	Speed      float64 `json:"speed"`
	ETA        float64 `json:"eta"`
}

// struct representing the response sent to the client
// as JSON-RPC result field
type ProcessResponse struct {
	Id       string              `json:"id"`
	Progress DownloadProgress    `json:"progress"`
	Info     common.DownloadInfo `json:"info"`
	Output   DownloadOutput      `json:"output"`
	Params   []string            `json:"params"`
}

// struct representing the current status of the memoryDB
// used for serializaton/persistence reasons
type Session struct {
	Processes []ProcessResponse `json:"processes"`
}

// struct representing the intent to stop a specific process
type AbortRequest struct {
	Id string `json:"id"`
}

// struct representing the intent to start a download
type DownloadRequest struct {
	Id     string
	URL    string   `json:"url"`
	Path   string   `json:"path"`
	Rename string   `json:"rename"`
	Params []string `json:"params"`
}

// struct representing request of creating a netscape cookies file
type SetCookiesRequest struct {
	Cookies string `json:"cookies"`
}

// represents a user defined collection of yt-dlp arguments
type CustomTemplate struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
