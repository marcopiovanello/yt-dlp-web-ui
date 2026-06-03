package formats

// Used to deser the formats in the -J output
type Metadata struct {
	Type      string     `json:"_type"`
	Formats   []Format   `json:"formats"`
	Best      Format     `json:"best"`
	Thumbnail string     `json:"thumbnail"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Entries   []Metadata `json:"entries"` // populated if url is playlist
}

func (m *Metadata) IsPlaylist() bool {
	return m.Type == "playlist"
}

// A skimmed yt-dlp format node
type Format struct {
	Format_id   string  `json:"format_id"`
	Format_note string  `json:"format_note"`
	FPS         float32 `json:"fps"`
	Resolution  string  `json:"resolution"`
	VCodec      string  `json:"vcodec"`
	ACodec      string  `json:"acodec"`
	Size        float64 `json:"filesize_approx"`
	Language    string  `json:"language"`
}
