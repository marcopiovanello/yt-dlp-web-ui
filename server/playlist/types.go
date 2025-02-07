package playlist

import "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"

type Metadata struct {
	Entries       []common.DownloadInfo `json:"entries"`
	Count         int                   `json:"playlist_count"`
	PlaylistTitle string                `json:"title"`
	Type          string                `json:"_type"`
}
