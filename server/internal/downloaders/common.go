package downloaders

import (
	"log/slog"
	"sync"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"
)

type DownloaderBase struct {
	Id        string
	URL       string
	Metadata  common.DownloadMetadata
	Pending   bool
	Completed bool
	mutex     sync.Mutex
}

func (d *DownloaderBase) FetchMetadata(fetcher func(url string) (*common.DownloadMetadata, error)) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	meta, err := fetcher(d.URL)
	if err != nil {
		slog.Error("failed to retrieve metadata", slog.Any("err", err))
		return
	}

	d.Metadata = *meta
}

func (d *DownloaderBase) SetPending(p bool) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Pending = p
}

func (d *DownloaderBase) Complete() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.Completed = true
}
