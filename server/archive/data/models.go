package data

import "time"

type ArchiveEntry struct {
	Id        string
	Title     string
	Path      string
	Thumbnail string
	Source    string
	Metadata  string
	CreatedAt time.Time
}
