package common

import "time"

// Used to deser the yt-dlp -J output
type DownloadMetadata struct {
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Thumbnail   string    `json:"thumbnail"`
	Resolution  string    `json:"resolution"`
	Size        int32     `json:"filesize_approx"`
	VCodec      string    `json:"vcodec"`
	ACodec      string    `json:"acodec"`
	Extension   string    `json:"ext"`
	OriginalURL string    `json:"original_url"`
	FileName    string    `json:"filename"`
	CreatedAt   time.Time `json:"created_at"`
}
