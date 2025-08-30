package playlist

import "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"

type Metadata struct {
	Entries       []common.DownloadMetadata `json:"entries"`
	Count         int                       `json:"playlist_count"`
	PlaylistTitle string                    `json:"title"`
	Type          string                    `json:"_type"`
}

func (m *Metadata) IsPlaylist() bool { return m.Type == "playlist" }
