package livestream

import "time"

type LiveStreamStatus = map[string]Status

type Status = struct {
	Status   int           `json:"status"`
	WaitTime time.Duration `json:"waitTime"`
	LiveDate time.Time     `json:"liveDate"`
}
