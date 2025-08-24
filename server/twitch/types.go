package twitch

import "time"

type StreamInfo struct {
	ID        string
	UserName  string
	Title     string
	GameName  string
	StartedAt time.Time
	IsLive    bool
}

type VodInfo struct {
	ID        string
	Title     string
	URL       string
	Duration  string
	CreatedAt time.Time
}
